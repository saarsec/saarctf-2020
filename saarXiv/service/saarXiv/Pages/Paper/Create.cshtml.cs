using System.ComponentModel.DataAnnotations;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Identity;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.RazorPages;
using saarXiv.Data;
using saarXiv.Models;

namespace saarXiv.Pages.Paper
{
    public class CreateModel : PageModel
    {
        private readonly ApplicationDbContext _context;
        private readonly UserManager<Author> _userManager;

        public CreateModel(UserManager<Author> userManager,
            ApplicationDbContext context)
        {
            _userManager = userManager;
            _context = context;
        }

        [BindProperty] public InputModel Input { get; set; }

        public class InputModel
        {
            [Display(Name = "Title")] public string Title { get; set; }

            [Display(Name = "Content")] public string Content { get; set; }
            
            [Display(Name = "This paper is currently under submission")] public bool UnderSubmission { get; set; }
        }

        public IActionResult OnGet()
        {
            Input = new InputModel
            {
                Title = "My glorious paper",
                Content = @"In this paper we explain yet another side-channel on Intel CPUs / prove P vs. NP / show the superios of tabs over spaces...",
                UnderSubmission = false
            };
            return Page();
        }

        // To protect from overposting attacks, please enable the specific properties you want to bind to, for
        // more details see https://aka.ms/RazorPagesCRUD.
        public async Task<IActionResult> OnPostAsync()
        {
            if (!ModelState.IsValid)
            {
                return Page();
            }

            var paper = new Models.Paper
            {
                Author = await _userManager.GetUserAsync(User),
                NeedsCompilation = true,
                Title = Input.Title,
                Content = Input.Content,
                UnderSubmission = Input.UnderSubmission
            };
            _context.Papers.Add(paper);
            await _context.SaveChangesAsync();
            TempData["StatusMessage"] = $"Successfully created paper {paper.ID}";
            return RedirectToPage("./Index");
        }
    }
}