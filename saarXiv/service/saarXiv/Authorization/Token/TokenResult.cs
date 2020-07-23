using System;
using System.Collections.Generic;
using System.Linq;

namespace saarXiv.Authorization.Token
{
    public class TokenResult
    {
        private static readonly TokenResult _success = new TokenResult()
        {
            Succeeded = true
        };

        public bool Succeeded { get; protected set; }

        public string Error { get; protected set; }
        

        public static TokenResult Success
        {
            get { return TokenResult._success; }
        }

        public static TokenResult Failed(string error)
        {
            return new TokenResult()
            {
                Succeeded = false,
                Error = error
            };
            
        }
        
    }
}