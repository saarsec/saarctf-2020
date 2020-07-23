from Crypto.PublicKey import RSA
from Crypto.Hash import SHA512
from Crypto.Signature import PKCS1_v1_5
from Crypto.Util.number import long_to_bytes, bytes_to_long


def sign(message: str, secret_key):
    signer = PKCS1_v1_5.new(secret_key)
    return signer.sign(SHA512.new(message.encode("utf-8")))


def verify(message: str, signature: bytes, public_key):
    verifier = PKCS1_v1_5.new(public_key)
    return verifier.verify(SHA512.new(message), signature)


def loadKey(filename, isPrivate):
    with open(filename, "rb") as f:
        key = RSA.importKey(f.read())
    if isPrivate:
        assert key.has_private(), "Keyfile does not contain a private key."
    else:
        assert not key.has_private(), "Keyfile should not contain a private key."
    return key


def encrypt(method, data, parameters):
    assert method in methods, "Unknown algorithm: " + method
    return methods[method](data, parameters)


def encryptPlain(data, parameters):
    # The DNPR made us very aware of our clients' desire for privacy. Hence, we don't even store their mail addresses
    # anywhere and can use this channel to update them with important news.
    assert parameters["method"] == "plain"
    return data


def encryptCaesar(data, parameters):
    # Maybe a bit dated, but some Saarländers are very fond of all the Roman stuff that was dug out here.
    # Fun fact: In some villages around here, people are afraid to even dig a hole in their garden. Should they by accident
    # unearth some Roman debris, things will get complicated because some monument conservation people will show up.
    assert parameters["method"] == "caesar"
    key = parameters["key"]
    return bytes([(x + key) % 256 for x in data])


def encryptSAAR(data, parameters):
    assert parameters["method"] == "saar"
    assert int(parameters["n"]) >> 2049 == 0, "n is too big"
    c = pow(bytes_to_long(data), int(parameters["e"]), int(parameters["n"]))
    return long_to_bytes(c)


def encryptOTP(data, parameters):
    # The OTP is the only scheme that's secure from an information theoretical point of view, so there's no reason not
    # to include it.
    assert parameters["method"] == "otp"
    from binascii import unhexlify
    encryptionKey = unhexlify(parameters["encryptionkey"])
    assert len(data) == len(encryptionKey), "Key and data must have the same length: %d != %d" % (
        len(data), len(encryptionKey))
    r = bytes([m ^ k for (m, k) in zip(data, encryptionKey)])
    return r


# We are very proud to present to you the latest encryption made in Saarland!
# It doesn't even use blockchain or AI technology at all. Just good, old-fashioned maths!
# The paper describing our "Super Cool but Horribly Weird ENKryption scheme" (SCHWENK scheme) is currently under submission
# to a top-tier security conference.

from Crypto.Util.number import inverse
class Schwenker1:
    schwenkingBook1 = dict()
    veryThickSchwenkingBook = dict()

    @staticmethod
    def schwenkingStuff(someNeatSchwenkerCount, sausagesAreMeatToo):
        if sausagesAreMeatToo in Schwenker1.schwenkingBook1:
            return Schwenker1.veryThickSchwenkingBook[sausagesAreMeatToo][someNeatSchwenkerCount], \
                   Schwenker1.schwenkingBook1[sausagesAreMeatToo]

        secretIngredient = len(bin(sausagesAreMeatToo)[2:])
        noVeggiesAllowed = Schwenker1.schwenkItBaby(sausagesAreMeatToo)
        tastySchwenker = len(bin(noVeggiesAllowed)[2:])
        aPieceOfTastySchwenker = tastySchwenker - 1
        charcoal = aPieceOfTastySchwenker
        anotherSchwenker = 30 * [0]
        fireplace = 0
        iAmHungry = noVeggiesAllowed >> (secretIngredient + 1)
        for countingSchwenker in range(30):
            fireplace += sausagesAreMeatToo * iAmHungry
            anotherSchwenker[countingSchwenker] = fireplace
            iAmHungry <<= aPieceOfTastySchwenker

        Schwenker1.veryThickSchwenkingBook[sausagesAreMeatToo] = anotherSchwenker
        Schwenker1.schwenkingBook1[sausagesAreMeatToo] = charcoal
        return anotherSchwenker[someNeatSchwenkerCount], charcoal

    @staticmethod
    def schwenkItBaby(deliciousSchwenker):
        SOME_MAGIC_SCHWENKING_NUMBER = 2
        THE_NUMBER_OF_BEERS_YOU_SHOULD_HAVE_WHILE_SCHWENKING = 100
        blazingFire = THE_NUMBER_OF_BEERS_YOU_SHOULD_HAVE_WHILE_SCHWENKING * deliciousSchwenker ** SOME_MAGIC_SCHWENKING_NUMBER
        isItDoneYet = len(bin(blazingFire)[2:])
        return 2 ** (isItDoneYet + 1)

    def __init__(self, pieceOfMeat, secretSchwenkerRecipe):
        self.theTastiestSchwenkerOfAll = pieceOfMeat
        self.fireplace = Schwenker1.schwenkItBaby(pieceOfMeat)
        self.anotherSchwenkerPlease = 0
        schwenkingFactor = 1
        for stuff in secretSchwenkerRecipe:
            self.anotherSchwenkerPlease += stuff * schwenkingFactor
            schwenkingFactor *= self.fireplace

    @staticmethod
    def gimmeASchwenker(where, isMy, schwenker):
        meat = schwenker * [0]
        iObviouslyKnowHowToCountSchwenkers = 0
        while where != 0:
            fire = where % isMy
            where = (where - fire) // isMy
            meat[iObviouslyKnowHowToCountSchwenkers] = fire
            iObviouslyKnowHowToCountSchwenkers += 1
        return Schwenker1(isMy, meat)

    @staticmethod
    def leetSchwenker():
        return Schwenker1(1337, [42])

    @staticmethod
    def convertItToASchwenker(notASchwenker, butThisIs):
        mealTime = Schwenker1.leetSchwenker()
        mealTime.anotherSchwenkerPlease = notASchwenker
        mealTime.theTastiestSchwenkerOfAll = butThisIs
        mealTime.fireplace = Schwenker1.schwenkItBaby(butThisIs)
        mealTime.turnTheSchwenkerAround()
        return mealTime

    def magicSchwenkerDuplication(self):
        schwenker = Schwenker1.leetSchwenker()
        schwenker.theTastiestSchwenkerOfAll = self.theTastiestSchwenkerOfAll
        schwenker.anotherSchwenkerPlease = self.anotherSchwenkerPlease
        schwenker.fireplace = self.fireplace
        return schwenker

    def turnTheSchwenkerAround(self):
        numberOfSchwenkers = Schwenker1.schwenkingStuff(0, self.theTastiestSchwenkerOfAll)[1]
        numberOfBeers = 16
        partyTime = self.anotherSchwenkerPlease

        yummy = (1 << (numberOfBeers - 1) * numberOfSchwenkers) - 1
        while (partyTime & yummy == partyTime and yummy != 0):
            numberOfBeers -= 1
            yummy >>= numberOfSchwenkers
        thinkingOfSoMuchSchwenkerstuffIsHard = 0
        schwenkersAreNoWeisswurst = numberOfBeers * numberOfSchwenkers
        yummy = (1 << numberOfSchwenkers) - 1
        eatingAPiece = 0
        while eatingAPiece <= schwenkersAreNoWeisswurst:
            someSchwenker = (partyTime & yummy) % self.theTastiestSchwenkerOfAll
            thinkingOfSoMuchSchwenkerstuffIsHard += someSchwenker << eatingAPiece
            partyTime >>= numberOfSchwenkers
            eatingAPiece += numberOfSchwenkers
        self.anotherSchwenkerPlease = thinkingOfSoMuchSchwenkerstuffIsHard

    def __add__(self, secondSchwenker):
        moreSchwenkerPlease = self.magicSchwenkerDuplication()
        moreSchwenkerPlease.anotherSchwenkerPlease = self.anotherSchwenkerPlease + secondSchwenker.n
        moreSchwenkerPlease.turnTheSchwenkerAround()
        return moreSchwenkerPlease

    def addSomeSpices(self, somePeopleCallSchwenkerASchaukelsteak):
        schwenkerCount = 16
        yetAnotherSchwenker = (1 << (schwenkerCount - 1) * somePeopleCallSchwenkerASchaukelsteak) - 1
        while (self.anotherSchwenkerPlease & yetAnotherSchwenker == self.anotherSchwenkerPlease):
            schwenkerCount -= 1
            yetAnotherSchwenker >>= somePeopleCallSchwenkerASchaukelsteak
        return schwenkerCount

    def cutMeAPieceOfSchwenker(self, evenMoreSchwenker):
        noCampingTripWithoutSchwenker = Schwenker1.schwenkingStuff(0, self.theTastiestSchwenkerOfAll)[1]

        aSadSchwenker = 0
        numberOfRemainingSchwenkers = self.anotherSchwenkerPlease
        emergencySchwenker = self.addSomeSpices(noCampingTripWithoutSchwenker)
        theresNeverTooMuchSchwenker = evenMoreSchwenker.addSomeSpices(noCampingTripWithoutSchwenker)
        numberOfRemainingSchwenkers += \
            Schwenker1.schwenkingStuff(emergencySchwenker - 1, self.theTastiestSchwenkerOfAll)[0]
        iAmHungry = (evenMoreSchwenker.anotherSchwenkerPlease >> (
                theresNeverTooMuchSchwenker - 1) * noCampingTripWithoutSchwenker) % self.theTastiestSchwenkerOfAll
        theBiggestSchwenkerOfAll = emergencySchwenker - 1
        while theBiggestSchwenkerOfAll >= theresNeverTooMuchSchwenker - 1:
            aSausageForAChange = (numberOfRemainingSchwenkers >> (
                    theBiggestSchwenkerOfAll * noCampingTripWithoutSchwenker)) % self.theTastiestSchwenkerOfAll
            sooooooManyyyyySchwenkers = (aSausageForAChange * inverse(iAmHungry,
                                                                      self.theTastiestSchwenkerOfAll)) % self.theTastiestSchwenkerOfAll
            aSimpleSchwenker = sooooooManyyyyySchwenkers << (
                    theBiggestSchwenkerOfAll - theresNeverTooMuchSchwenker + 1) * noCampingTripWithoutSchwenker
            aSadSchwenker += aSimpleSchwenker
            numberOfRemainingSchwenkers -= evenMoreSchwenker.anotherSchwenkerPlease * aSimpleSchwenker
            theBiggestSchwenkerOfAll -= 1

        return Schwenker1.convertItToASchwenker(aSadSchwenker,
                                                self.theTastiestSchwenkerOfAll), Schwenker1.convertItToASchwenker(
            numberOfRemainingSchwenkers, self.theTastiestSchwenkerOfAll)

    def iAmNotHungry(self):
        theSmallesSchwenkerOfAll = self.magicSchwenkerDuplication()
        theSmallesSchwenkerOfAll.anotherSchwenkerPlease = 0
        return theSmallesSchwenkerOfAll

    def iAmALittleHungry(self):
        oneSchwenker = self.magicSchwenkerDuplication()
        oneSchwenker.anotherSchwenkerPlease = 1
        return oneSchwenker

    def gimmeSomeSchwenker(self, schwenkingIsFun):
        if self.anotherSchwenkerPlease == 0:
            return schwenkingIsFun, self.iAmNotHungry(), self.iAmALittleHungry()

        aPieceOfSchwenker, maybeAWholeSchwenker = schwenkingIsFun.cutMeAPieceOfSchwenker(self)
        schw, en, ker = maybeAWholeSchwenker.gimmeSomeSchwenker(self)
        return (schw, Schwenker1.convertItToASchwenker(
            ker.anotherSchwenkerPlease + Schwenker1.schwenkingStuff(16, self.theTastiestSchwenkerOfAll)[
                0] - aPieceOfSchwenker.anotherSchwenkerPlease * en.anotherSchwenkerPlease,
            self.theTastiestSchwenkerOfAll), en)


class MysteriousSchwenker:
    schwenkStuff = dict()

    @staticmethod
    def prepareSchwenker(schwenker1, schwenker2):
        if (schwenker1, schwenker2) in MysteriousSchwenker.schwenkStuff:
            return MysteriousSchwenker.schwenkStuff[(schwenker1, schwenker2)]

        schwenker3 = Schwenker1.schwenkItBaby(schwenker1)
        schwenker4 = len(bin(schwenker3)[2:])
        schwenker5 = len(bin(schwenker1)[2:])
        schwenker6 = len(bin(schwenker2)[2:]) // (schwenker4 - 1) + 1
        schwenker7 = schwenker4 - 1

        schwenker8 = 30 * [0]
        schwenker10 = 0
        schwenker11 = schwenker3 >> (schwenker5 + 1)
        for schwenker12 in range(30):
            schwenker10 += schwenker1 * schwenker11
            schwenker8[schwenker12] = schwenker10
            schwenker11 <<= schwenker7

        schwenker13 = (schwenker6, schwenker7, schwenker8)
        MysteriousSchwenker.schwenkStuff[(schwenker1, schwenker2)] = schwenker13
        return schwenker13

    def __init__(self, publicSchwenker, secretSchwenker, signingSchwenker, verificationSchwenker):
        self.publicSchwenker = publicSchwenker
        self.secretSchwenker = secretSchwenker
        self.signingSchwenker = signingSchwenker.anotherSchwenkerPlease
        self.verificationSchwenker = verificationSchwenker.anotherSchwenkerPlease

    @staticmethod
    def eatASchwenker(moreSchwenker, we, want, schwenker):
        pieceOfSchwenker = Schwenker1(we, schwenker)
        yummy = Schwenker1.gimmeASchwenker(moreSchwenker, we, want)
        return MysteriousSchwenker(we, want, pieceOfSchwenker, yummy)

    def eatAnotherSchwenker(self):
        schwenkerSmellsDelicious = Schwenker1.schwenkItBaby(self.publicSchwenker)
        schwenkerStack = []
        schwenkerCount = self.verificationSchwenker
        schwenkerIterator = 0
        while schwenkerCount != 0:
            schwenkerCount, anotherSchwenkerCount = divmod(schwenkerCount, schwenkerSmellsDelicious)
            schwenkerStack += [anotherSchwenkerCount]
            schwenkerIterator += 1
        return sum(deliciousSchwenker * (self.publicSchwenker ** veryDeliciousSchwenker) for
                   veryDeliciousSchwenker, deliciousSchwenker in enumerate(schwenkerStack))

    @staticmethod
    def gimmeTheBestSchwenkerOfAll():
        epicSchwenker = Schwenker1.leetSchwenker()
        return MysteriousSchwenker(1337, 42, epicSchwenker, epicSchwenker)

    def magicSchwenkerDuplication(self):
        imitationSchwenker = MysteriousSchwenker.gimmeTheBestSchwenkerOfAll()
        imitationSchwenker.publicSchwenker = self.publicSchwenker
        imitationSchwenker.secretSchwenker = self.secretSchwenker
        imitationSchwenker.signingSchwenker = self.signingSchwenker
        imitationSchwenker.verificationSchwenker = self.verificationSchwenker
        return imitationSchwenker

    def cutTheSchwenkerInHalf(self):
        s, c, h = MysteriousSchwenker.prepareSchwenker(self.publicSchwenker, self.signingSchwenker)

        w = self.verificationSchwenker
        e = 2 * (s - 1) - 1
        n = (1 << (e - 1) * c) - 1
        while (w & n == w):
            e -= 1
            n >>= c
        w += h[e - 1]
        k = e - 1
        r = (e - s) * c

        while r >= 0:
            S = (w >> (k * c)) % self.publicSchwenker
            w -= S * self.signingSchwenker << r
            w &= n
            k -= 1
            r -= c
            n >>= c

        if s <= e:
            C = s - 2
        else:
            C = e - 1

        H = 0
        W = C * c
        n = (1 << c) - 1
        E = 0
        while E <= W:
            N = (w & n) % self.publicSchwenker
            H += N << E
            w >>= c
            E += c
        self.verificationSchwenker = H

    def __add__(self, theOtherSchwenker):
        theOtherOtherSchwenker = self.magicSchwenkerDuplication()
        theOtherOtherSchwenker.verificationSchwenker = self.verificationSchwenker + theOtherSchwenker.verificationSchwenker
        theOtherOtherSchwenker.cutTheSchwenkerInHalf()
        return theOtherOtherSchwenker

    def __mul__(self, yetAnotherSchwenker):
        yetAnotherAnotherSchwenker = self.magicSchwenkerDuplication()
        yetAnotherAnotherSchwenker.verificationSchwenker = self.verificationSchwenker * yetAnotherSchwenker.verificationSchwenker
        yetAnotherAnotherSchwenker.cutTheSchwenkerInHalf()
        return yetAnotherAnotherSchwenker

    def roastTheSchwenker(self):
        moreSchwenker, evenMoreSchwenker, theMostSchwenker = Schwenker1.convertItToASchwenker(
            self.verificationSchwenker, self.publicSchwenker).gimmeSomeSchwenker(
            Schwenker1.convertItToASchwenker(self.signingSchwenker, self.publicSchwenker))
        theMostestSchwenker = Schwenker1.schwenkingStuff(0, self.publicSchwenker)[1]
        theMostestSchwenkerInTheWorld = moreSchwenker.addSomeSpices(theMostestSchwenker)
        theMostestSchwenkerInTheUniverse = (moreSchwenker.anotherSchwenkerPlease >> (
                theMostestSchwenkerInTheWorld - 1) * theMostestSchwenker) % self.publicSchwenker
        theMostestSchwenkerInTheMultiverse = self.magicSchwenkerDuplication()
        theMostestSchwenkerInTheMultiverse.verificationSchwenker = evenMoreSchwenker.anotherSchwenkerPlease * inverse(
            theMostestSchwenkerInTheUniverse, self.publicSchwenker)
        theMostestSchwenkerInTheMultiverse.cutTheSchwenkerInHalf()
        return theMostestSchwenkerInTheMultiverse

    def oneSchwenkerPlease(self):
        lonelySchwenker = self.magicSchwenkerDuplication()
        lonelySchwenker.verificationSchwenker = 1
        return lonelySchwenker

    def __pow__(self, powerfulSchwenker):
        everGrowingSchwenker = self.oneSchwenkerPlease()
        for tinyPieceOfSchwenker in bin(powerfulSchwenker)[2:]:
            everGrowingSchwenker = everGrowingSchwenker * everGrowingSchwenker
            if tinyPieceOfSchwenker == '1':
                everGrowingSchwenker = everGrowingSchwenker * self
        return everGrowingSchwenker


class Schwencryptor:
    MAGIC_SCHWENKING_CONSTANT = 32

    @staticmethod
    def S(schwenker1, schwenker2):
        return schwenker1 + schwenker2

    @staticmethod
    def C(halfOfASchwenker, theOtherHalfOfASchwenker):
        return halfOfASchwenker * theOtherHalfOfASchwenker.roastTheSchwenker(), Schwencryptor.S(
            theOtherHalfOfASchwenker, halfOfASchwenker)

    @staticmethod
    def H(notVeryPowerfulSchwenker, alsoNotAVeryPowerfulSchwenker, maybeAPowerfulSchwenker):
        return alsoNotAVeryPowerfulSchwenker ** (
                notVeryPowerfulSchwenker.eatAnotherSchwenker() ^ maybeAPowerfulSchwenker)

    @staticmethod
    def hideMySchwenker(rawSchwenker, schwenkerCookbook):
        from hashlib import sha512
        from Crypto.Util.number import long_to_bytes
        assert len(rawSchwenker) >= 64
        someSchwenker = int(rawSchwenker[:Schwencryptor.MAGIC_SCHWENKING_CONSTANT].hex(), 16)
        areYouReadyForSomeSchwenker = int(
            rawSchwenker[Schwencryptor.MAGIC_SCHWENKING_CONSTANT:2 * Schwencryptor.MAGIC_SCHWENKING_CONSTANT].hex(), 16)
        rawSchwenker = rawSchwenker[2 * Schwencryptor.MAGIC_SCHWENKING_CONSTANT:]
        rawSchwenker = Schwencryptor.putTheSchwenkerOnAPlate(rawSchwenker)

        assert len(rawSchwenker) % Schwencryptor.MAGIC_SCHWENKING_CONSTANT == 0
        thatsANiceNumberOfSchwenkers = len(rawSchwenker) // Schwencryptor.MAGIC_SCHWENKING_CONSTANT
        deliciousSchwenker = b""
        secretSchwenkerIngredient = schwenkerCookbook["schwenkerid"]
        notSoSecretSchwenkerIngredient = schwenkerCookbook["schwenkingoptions"]
        for countYourSchwenkers in range(thatsANiceNumberOfSchwenkers):
            partialSchwenker = rawSchwenker[
                               countYourSchwenkers * Schwencryptor.MAGIC_SCHWENKING_CONSTANT:(countYourSchwenkers + 1) * Schwencryptor.MAGIC_SCHWENKING_CONSTANT]
            schwenker, someOtherSchwenker, shcwenkre = Schwencryptor.getStuff(
                notSoSecretSchwenkerIngredient[countYourSchwenkers], someSchwenker, areYouReadyForSomeSchwenker,
                partialSchwenker)
            yummySchwenker = Schwencryptor.H(schwenker, shcwenkre, secretSchwenkerIngredient[countYourSchwenkers])
            deliciousSchwenker += long_to_bytes(yummySchwenker.eatAnotherSchwenker())
            someSchwenker = schwenker.eatAnotherSchwenker()
            areYouReadyForSomeSchwenker = someOtherSchwenker.eatAnotherSchwenker()

        noSchwenkerButMacAndCheese = sha512(deliciousSchwenker).digest()
        partialSchwenker = noSchwenkerButMacAndCheese[:Schwencryptor.MAGIC_SCHWENKING_CONSTANT]
        schwenker, someOtherSchwenker, shcwenkre = Schwencryptor.getStuff(notSoSecretSchwenkerIngredient[-2],
                                                                          someSchwenker, areYouReadyForSomeSchwenker,
                                                                          partialSchwenker)
        yummySchwenker = Schwencryptor.H(shcwenkre, schwenker, secretSchwenkerIngredient[-2])
        deliciousSchwenker += long_to_bytes(yummySchwenker.eatAnotherSchwenker())
        someSchwenker = schwenker.eatAnotherSchwenker()
        areYouReadyForSomeSchwenker = someOtherSchwenker.eatAnotherSchwenker()
        partialSchwenker = noSchwenkerButMacAndCheese[Schwencryptor.MAGIC_SCHWENKING_CONSTANT:]
        schwenker, someOtherSchwenker, shcwenkre = Schwencryptor.getStuff(notSoSecretSchwenkerIngredient[-1],
                                                                          someSchwenker, areYouReadyForSomeSchwenker,
                                                                          partialSchwenker)
        yummySchwenker = Schwencryptor.H(shcwenkre, schwenker, secretSchwenkerIngredient[-1])
        deliciousSchwenker += long_to_bytes(yummySchwenker.eatAnotherSchwenker())

        return deliciousSchwenker

    @staticmethod
    def getStuff(schwenker1, schwenker2, schwenker3, schwenker4):
        from Crypto.Util.number import bytes_to_long
        isMySchwenkerReadyYet = schwenker1[0]
        cmonIAmHungry = len(schwenker1) - 2
        hurryUp = schwenker1[1:]
        schwenkFaster = MysteriousSchwenker.eatASchwenker(schwenker2, isMySchwenkerReadyYet, cmonIAmHungry, hurryUp)
        reallyIMeanIt = MysteriousSchwenker.eatASchwenker(schwenker3, isMySchwenkerReadyYet, cmonIAmHungry, hurryUp)
        now = MysteriousSchwenker.eatASchwenker(bytes_to_long(schwenker4), isMySchwenkerReadyYet, cmonIAmHungry,
                                                hurryUp)
        atLast, iCanEatMyDeliciousSchwenker = Schwencryptor.C(schwenkFaster, reallyIMeanIt)
        return atLast, iCanEatMyDeliciousSchwenker, now

    @staticmethod
    def putTheSchwenkerOnAPlate(couldItBeAnotherSchwenker):
        if len(couldItBeAnotherSchwenker) % Schwencryptor.MAGIC_SCHWENKING_CONSTANT == 0:
            couldItBeAnotherSchwenker += Schwencryptor.MAGIC_SCHWENKING_CONSTANT * chr(
                Schwencryptor.MAGIC_SCHWENKING_CONSTANT).encode("utf-8")
        else:
            couldItBeAnotherSchwenker += (Schwencryptor.MAGIC_SCHWENKING_CONSTANT - len(
                couldItBeAnotherSchwenker) % Schwencryptor.MAGIC_SCHWENKING_CONSTANT) * chr(
                Schwencryptor.MAGIC_SCHWENKING_CONSTANT - len(
                    couldItBeAnotherSchwenker) % Schwencryptor.MAGIC_SCHWENKING_CONSTANT).encode("utf-8")
        return couldItBeAnotherSchwenker


def schwencrypt(publicSchwenker, secretSchwenker):
    return Schwencryptor.hideMySchwenker(publicSchwenker, secretSchwenker)


methods = dict({
    "plain": encryptPlain,
    "caesar": encryptCaesar,
    "saar": encryptSAAR,
    "otp": encryptOTP,
    "schwenk": schwencrypt
})
