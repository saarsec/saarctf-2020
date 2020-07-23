using Microsoft.AspNetCore.Authorization.Infrastructure;

namespace saarXiv.Authorization.Operations
{
    public class PaperOperations
    {
        public static OperationAuthorizationRequirement Download =
            new OperationAuthorizationRequirement {Name = nameof(Download)};

        public static OperationAuthorizationRequirement Share =
            new OperationAuthorizationRequirement {Name = nameof(Share)};
        
        public static OperationAuthorizationRequirement Edit =
            new OperationAuthorizationRequirement {Name = nameof(Edit)};

        public static OperationAuthorizationRequirement Delete =
            new OperationAuthorizationRequirement {Name = nameof(Delete)};
    }
}