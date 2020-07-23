using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Identity;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.RazorPages;
using Microsoft.EntityFrameworkCore;
using saarXiv.Data;
using saarXiv.Models;

namespace saarXiv.Pages.Paper
{
    public class IndexModel : PageModel
    {
        private readonly UserManager<Author> _userManager;
        private readonly ApplicationDbContext _context;

        public IndexModel(UserManager<Author> userManager,
            ApplicationDbContext context)
        {
            _userManager = userManager;
            _context = context;
        }

        public IList<Models.Paper> Paper { get; set; }

        [TempData] public string StatusMessage { get; set; }

        public async Task<IActionResult> OnGetAsync()
        {
            var user = await _userManager.GetUserAsync(User);
            if (user == null)
            {
                return NotFound($"Unable to load user with ID '{_userManager.GetUserId(User)}'.");
            }

            Paper = await _context.Papers.Where(paper => paper.Author.Equals(user)).ToListAsync();
            return Page();
        }
    }
}