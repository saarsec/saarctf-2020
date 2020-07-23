using System;
using System.Linq;
using System.Security.Cryptography;
using Microsoft.AspNetCore.Identity;

namespace saarXiv.Models
{
    public class Author : IdentityUser
    {
        [ProtectedPersonalData]
        public string Firstname { get; set; }

        [ProtectedPersonalData]
        public string Lastname { get; set; }

        [ProtectedPersonalData]
        public string Key { get; set; }

        public string Initial
        {
            get
            {
                if (Firstname.Length > 0)
                {
                    return Firstname.Substring(0, 1);
                }
                else
                {
                    return "";
                }
            }
        }
        
        public ECDsa ECDsaKey
        {
            get
            {
                var ecdsa = ECDsa.Create();
                byte[] keyBytes = Convert.FromBase64String(Key);
                int bytesRead;
                ecdsa.ImportPkcs8PrivateKey(keyBytes, out bytesRead);
                return ecdsa;
            }
        }

        public byte[] PublicKey
        {
            get
            {
                var publicParameters = ECDsaKey.ExportParameters(false);
                var publicKey = new byte[publicParameters.Q.X.Length + publicParameters.Q.Y.Length + 1];
                publicKey[0] = 0x04;
                Array.Copy(publicParameters.Q.X, 0, publicKey, 1, publicParameters.Q.X.Length);
                Array.Copy(publicParameters.Q.Y, 0, publicKey, publicParameters.Q.X.Length+1, publicParameters.Q.Y.Length);
                return publicKey;
            }
        }
    }
}