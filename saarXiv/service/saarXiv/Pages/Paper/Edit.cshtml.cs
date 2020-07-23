using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.RazorPages;
using Microsoft.EntityFrameworkCore;
using saarXiv.Authorization.Operations;
using saarXiv.Data;

namespace saarXiv.Pages.Paper
{
    public class EditModel : PageModel
    {
        private readonly ApplicationDbContext _context;
        private readonly IAuthorizationService _authorizationService;

        public EditModel(IAuthorizationService authorizationService, ApplicationDbContext context)
        {
            _authorizationService = authorizationService;
            _context = context;
        }

        [BindProperty] public InputModel Input { get; set; }

        public class InputModel
        {
            [Display(Name = "Title")] public string Title { get; set; }

            [Display(Name = "Content")] public string Content { get; set; }
            
            [Display(Name = "This paper is currently under submission")] public bool UnderSubmission { get; set; }
        }

        public async Task<IActionResult> OnGetAsync(int? id)
        {
            if (id == null)
            {
                return NotFound();
            }

            var paper = await _context.Papers.FirstOrDefaultAsync(m => m.ID == id);

            if (paper == null)
            {
                return NotFound();
            }

            var authorizationResult = await _authorizationService
                .AuthorizeAsync(User, paper, PaperOperations.Edit);

            Input = new InputModel
            {
                Title = paper.Title, Content = paper.Content, UnderSubmission = paper.UnderSubmission
            };

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

        // To protect from overposting attacks, please enable the specific properties you want to bind to, for
        // more details see https://aka.ms/RazorPagesCRUD.
        public async Task<IActionResult> OnPostAsync(int? id)
        {
            if (!ModelState.IsValid)
            {
                return Page();
            }

            if (id == null)
            {
                return NotFound();
            }

            var paper = await _context.Papers.FirstOrDefaultAsync(m => m.ID == id);

            if (paper == null)
            {
                return NotFound();
            }

            var authorizationResult = await _authorizationService
                .AuthorizeAsync(User, paper, PaperOperations.Edit);

            paper.Title = Input.Title;
            paper.Content = Input.Content;
            paper.UnderSubmission = Input.UnderSubmission;
            paper.NeedsCompilation = true;
            _context.Attach(paper).State = EntityState.Modified;

            try
            {
                await _context.SaveChangesAsync();
            }
            catch (DbUpdateConcurrencyException)
            {
                if (!PaperExists(paper.ID))
                {
                    return NotFound();
                }
                else
                {
                    throw;
                }
            }

            return RedirectToPage("./Index");
        }

        private bool PaperExists(int id)
        {
            return _context.Papers.Any(e => e.ID == id);
        }
    }
}