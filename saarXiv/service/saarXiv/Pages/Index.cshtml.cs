using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Identity;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.RazorPages;
using Microsoft.EntityFrameworkCore;
using saarXiv.Data;
using saarXiv.Models;

namespace saarXiv.Pages
{
    public class IndexModel : PageModel
    {
        private readonly ApplicationDbContext _context;

        public IndexModel(UserManager<Author> userManager,
            ApplicationDbContext context)
        {
            _context = context;
        }

        public IList<Models.Paper> Papers { get; set; }

        public async Task<IActionResult> OnGetAsync()
        {
            Papers = await _context.Papers.Include(paper => paper.Author).OrderByDescending(paper => paper.ID).Take(16)
                .ToListAsync();
            return Page();
        }
    }
}