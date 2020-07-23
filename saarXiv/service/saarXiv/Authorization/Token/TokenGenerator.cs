using System;
using System.Security.Cryptography;
using System.Security.Cryptography.X509Certificates;
using Microsoft.AspNetCore.Identity;
using saarXiv.Models;

namespace saarXiv.Authorization.Token
{
    public class TokenGenerator
    {

        public string GenerateTokenAsync(Paper paper)
        {
            var request = new CertificateRequest(new X500DistinguishedName($"CN={paper.SaarXivID}"), paper.Author.ECDsaKey, HashAlgorithmName.SHA256);
            var cert = request.CreateSelfSigned(DateTimeOffset.Now, DateTimeOffset.Now.AddDays(1));
            return Convert.ToBase64String(cert.GetRawCertData());
        }
    }
}