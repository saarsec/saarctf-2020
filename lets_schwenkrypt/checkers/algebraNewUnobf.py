from Crypto.Util.number import inverse


class ZpXElement:
    lastDigitBits = dict()
    consts = dict()

    @staticmethod
    def getFieldParams(digits, p):
        if p in ZpXElement.lastDigitBits:
            return ZpXElement.consts[p][digits], ZpXElement.lastDigitBits[p]

        pBits = len(bin(p)[2:])
        base = ZpXElement.computeBase(p)
        baseBits = len(bin(base)[2:])
        digitBits = baseBits - 1
        lastDigitBits = digitBits
        consts = 30 * [0]
        const = 0
        nextMult = base >> (pBits + 1)
        for i in range(30):
            const += p * nextMult
            consts[i] = const
            nextMult <<= digitBits

        ZpXElement.consts[p] = consts
        ZpXElement.lastDigitBits[p] = lastDigitBits
        return consts[digits], lastDigitBits

    @staticmethod
    def computeBase(p):
        SPACE_POWER_MULTIPLIER = 2
        SPACE_MULTIPLIER = 100
        base = SPACE_MULTIPLIER * p ** SPACE_POWER_MULTIPLIER
        bits = len(bin(base)[2:])
        return 2 ** (bits + 1)

    def __init__(self, p, coeffs):
        self.p = p
        self.base = ZpXElement.computeBase(p)
        self.n = 0
        curMult = 1
        for c in coeffs:
            self.n += c * curMult
            curMult *= self.base

    @staticmethod
    def fromDenseNumber(n, b, p):
        coeffs = p * [0]
        i = 0
        while n != 0:
            c = n % b
            n = (n - c) // b
            coeffs[i] = c
            i += 1
        return ZpXElement(b, coeffs)

    @staticmethod
    def dummy():
        return ZpXElement(1337, [42])

    @staticmethod
    def fromNumber(n, p):
        d = ZpXElement.dummy()
        d.n = n
        d.p = p
        d.base = ZpXElement.computeBase(p)
        d.applyModulus()
        return d

    def clone(self):
        d = ZpXElement.dummy()
        d.p = self.p
        d.n = self.n
        d.base = self.base
        return d

    def applyModulus(self):
        digitBits = ZpXElement.getFieldParams(0, self.p)[1]
        nDigits = 16
        n = self.n

        mask = (1 << (nDigits - 1) * digitBits) - 1
        while (n & mask == n and mask != 0):
            nDigits -= 1
            mask >>= digitBits
        res = 0
        target = nDigits * digitBits
        mask = (1 << digitBits) - 1
        shift = 0
        while shift <= target:
            coeff = (n & mask) % self.p
            res += coeff << shift
            n >>= digitBits
            shift += digitBits
        self.n = res

    def __add__(self, rhs):
        res = self.clone()
        res.n = self.n + rhs.n
        res.applyModulus()
        return res

    def computeDigits(self, digitBits):
        nDigits = 16
        mask = (1 << (nDigits - 1) * digitBits) - 1
        while (self.n & mask == self.n):
            nDigits -= 1
            mask >>= digitBits
        return nDigits

    def polyLongDiv(self, rhs):
        digitBits = ZpXElement.getFieldParams(0, self.p)[1]

        quotient = 0
        remainder = self.n
        remainderDigits = self.computeDigits(digitBits)
        bDigits = rhs.computeDigits(digitBits)
        remainder += ZpXElement.getFieldParams(remainderDigits - 1, self.p)[0]
        divisorLeadingCoeff = (rhs.n >> (bDigits - 1) * digitBits) % self.p
        digit = remainderDigits - 1
        while digit >= bDigits - 1:
            remainderLeadingCoeff = (remainder >> (digit * digitBits)) % self.p
            currentQuotMult = (remainderLeadingCoeff * inverse(divisorLeadingCoeff, self.p)) % self.p
            currentQuotMonomial = currentQuotMult << (digit - bDigits + 1) * digitBits
            quotient += currentQuotMonomial
            remainder -= rhs.n * currentQuotMonomial
            digit -= 1

        return ZpXElement.fromNumber(quotient, self.p), ZpXElement.fromNumber(remainder, self.p)

    def zero(self):
        z = self.clone()
        z.n = 0
        return z

    def one(self):
        o = self.clone()
        o.n = 1
        return o

    def egcd(self, rhs):
        if self.n == 0:
            return rhs, self.zero(), self.one()

        quot, rem = rhs.polyLongDiv(self)
        g, y, x = rem.egcd(self)
        return (g, ZpXElement.fromNumber(x.n + ZpXElement.getFieldParams(16, self.p)[0] - quot.n * y.n, self.p), y)


class GFElement:
    precomputed = dict()

    @staticmethod
    def precompute(p, nModulus):
        if (p, nModulus) in GFElement.precomputed:
            return GFElement.precomputed[(p, nModulus)]

        base = ZpXElement.computeBase(p)
        baseBits = len(bin(base)[2:])
        pBits = len(bin(p)[2:])
        modDigits = len(bin(nModulus)[2:]) // (baseBits - 1) + 1
        digitBits = baseBits - 1

        constants = 30 * [0]
        const = 0
        nextMult = base >> (pBits + 1)
        for i in range(30):
            const += p * nextMult
            constants[i] = const
            nextMult <<= digitBits

        precomputed = (modDigits, digitBits, constants)
        GFElement.precomputed[(p, nModulus)] = precomputed
        return precomputed

    def __init__(self, p: int, power: int, modulus: ZpXElement, elem: ZpXElement):
        self.p = p
        self.power = power
        self.modulus = modulus.n
        self.n = elem.n

    @staticmethod
    def fromDenseNumber(n, base, power, modulusCoeffs):
        assert base ** power >= n, "Not enough space: %d < %d" % (base ** power, n)
        assert n >= 0
        modulus = ZpXElement(base, modulusCoeffs)
        elem = ZpXElement.fromDenseNumber(n, base, power)
        return GFElement(base, power, modulus, elem)

    def toDenseNumber(self):
        base = ZpXElement.computeBase(self.p)
        coeffs = []
        n = self.n
        i = 0
        while n != 0:
            n, c = divmod(n, base)
            coeffs += [c]
            i += 1
        return sum(y * (self.p ** x) for x, y in enumerate(coeffs))

    @staticmethod
    def dummy():
        d1 = ZpXElement.dummy()
        return GFElement(1337, 42, d1, d1)

    def clone(self):
        d = GFElement.dummy()
        d.p = self.p
        d.power = self.power
        d.modulus = self.modulus
        d.n = self.n
        return d

    def applyModulus(self):
        modDigits, digitBits, constants = GFElement.precompute(self.p, self.modulus)

        n = self.n
        nDigits = 2 * (modDigits - 1) - 1
        mask = (1 << (nDigits - 1) * digitBits) - 1
        while (n & mask == n):
            nDigits -= 1
            mask >>= digitBits
        n += constants[nDigits - 1]
        digit = nDigits - 1
        modShift = (nDigits - modDigits) * digitBits

        while modShift >= 0:
            modMult = (n >> (digit * digitBits)) % self.p
            n -= modMult * self.modulus << modShift
            n &= mask
            digit -= 1
            modShift -= digitBits
            mask >>= digitBits

        if modDigits <= nDigits:
            numDigits = modDigits - 2
        else:
            numDigits = nDigits - 1

        res = 0
        target = numDigits * digitBits
        mask = (1 << digitBits) - 1
        shift = 0
        while shift <= target:
            coeff = (n & mask) % self.p
            res += coeff << shift
            n >>= digitBits
            shift += digitBits
        self.n = res

    def __add__(self, rhs):
        res = self.clone()
        res.n = self.n + rhs.n
        res.applyModulus()
        return res

    def __mul__(self, rhs):
        res = self.clone()
        res.n = self.n * rhs.n
        res.applyModulus()
        return res

    def inverse(self):
        g, x, y = ZpXElement.fromNumber(self.n, self.p).egcd(ZpXElement.fromNumber(self.modulus, self.p))
        digitBits = ZpXElement.getFieldParams(0, self.p)[1]
        gDigits = g.computeDigits(digitBits)
        gLeadingCoeff = (g.n >> (gDigits - 1) * digitBits) % self.p
        res = self.clone()
        res.n = x.n * inverse(gLeadingCoeff, self.p)
        res.applyModulus()
        return res

    def one(self):
        o = self.clone()
        o.n = 1
        return o

    def __pow__(self, exp):
        assert exp >= 0
        res = self.one()
        for bit in bin(exp)[2:]:
            res = res * res
            if bit == '1':
                res = res * self
        return res
