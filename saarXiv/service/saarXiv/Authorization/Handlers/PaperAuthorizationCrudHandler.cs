using System.Threading.Tasks;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Authorization.Infrastructure;
using Microsoft.AspNetCore.Identity;
using saarXiv.Authorization.Operations;
using saarXiv.Models;

namespace saarXiv.Authorization.Handlers
{
    public class PaperAuthorizationCrudHandler : AuthorizationHandler<OperationAuthorizationRequirement, Paper>
    {
        private readonly UserManager<Author> _userManager;

        public PaperAuthorizationCrudHandler(UserManager<Author> userManager)
        {
            _userManager = userManager;
        }

        protected override async Task HandleRequirementAsync(AuthorizationHandlerContext context,
            OperationAuthorizationRequirement requirement,
            Paper paper)
        {
            if (requirement == PaperOperations.Download && !paper.UnderSubmission)
            {
                context.Succeed(requirement);
            }

            else if (requirement == PaperOperations.Download || requirement == PaperOperations.Edit ||
                     requirement == PaperOperations.Delete || requirement == PaperOperations.Share)
            {
                var userIsAuthor = await CurrentUserIsAuthor(context, paper);
                if (userIsAuthor)
                {
                    context.Succeed(requirement);
                }
            }
        }

        private async Task<bool> CurrentUserIsAuthor(AuthorizationHandlerContext context, Paper paper)
        {
            var user = await _userManager.GetUserAsync(context.User);
            if (user == null)
            {
                return false;
            }

            return paper.Author == user;
        }
    }
}