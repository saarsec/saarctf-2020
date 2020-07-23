using System.Threading.Tasks;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.RazorPages;
using Microsoft.EntityFrameworkCore;
using saarXiv.Authorization.Operations;
using saarXiv.Data;

namespace saarXiv.Pages.Paper
{
    public class DeleteModel : PageModel
    {
        private readonly ApplicationDbContext _context;
        private IAuthorizationService _authorizationService;

        public DeleteModel(IAuthorizationService authorizationService,
            ApplicationDbContext context)
        {
            _authorizationService = authorizationService;
            _context = context;
        }

        [BindProperty] public Models.Paper Paper { get; set; }

        public async Task<IActionResult> OnGetAsync(int? id)
        {
            if (id == null)
            {
                return NotFound();
            }

            Paper = await _context.Papers.FirstOrDefaultAsync(m => m.ID == id);

            if (Paper == null)
            {
                return NotFound();
            }

            var authorizationResult = await _authorizationService
                .AuthorizeAsync(User, Paper, PaperOperations.Delete);

            if (authorizationResult.Succeeded)
            {
                return Page();
            }
            else if (User.Identity.IsAuthenticated)
            {
                return new ForbidResult();
            }
            else
            {
                return new ChallengeResult();
            }
        }

        public async Task<IActionResult> OnPostAsync(int? id)
        {
            if (id == null)
            {
                return NotFound();
            }

            Paper = await _context.Papers.FindAsync(id);

            if (Paper == null)
            {
                return NotFound();
            }

            var authorizationResult = await _authorizationService
                .AuthorizeAsync(User, Paper, PaperOperations.Delete);

            if (authorizationResult.Succeeded)
            {
                _context.Papers.Remove(Paper);
                await _context.SaveChangesAsync();

                return RedirectToPage("./Index");
            }
            else if (User.Identity.IsAuthenticated)
            {
                return new ForbidResult();
            }
            else
            {
                return new ChallengeResult();
            }
        }
    }
}