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
    public class AllModel : PageModel
    {
        private readonly ApplicationDbContext _context;

        public AllModel(UserManager<Author> userManager,
            ApplicationDbContext context)
        {
            _context = context;
        }

        public IList<Models.Paper> Papers { get; set; }

        public int PageIndex { get; set; }

        public int PageSize { get; set; }

        public bool HasPreviousPage
        {
            get { return PageIndex > 0; }
        }

        public bool HasNextPage
        {
            get { return Papers.Count() == PageSize; }
        }

        public async Task<IActionResult> OnGetAsync(int? p, int? q)
        {
            PageIndex = p ?? 0;
            PageSize = q ?? 16;
            int skip = PageIndex * PageSize;
            Papers = await _context.Papers.Include(paper => paper.Author).OrderByDescending(paper => paper.ID)
                .Skip(skip).Take(PageSize).ToListAsync();
            return Page();
        }
    }
}