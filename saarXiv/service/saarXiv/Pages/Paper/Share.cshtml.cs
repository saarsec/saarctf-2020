using System.IO;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.RazorPages;
using Microsoft.EntityFrameworkCore;
using saarXiv.Authorization.Operations;
using saarXiv.Authorization.Token;
using saarXiv.Controllers;
using saarXiv.Data;

namespace saarXiv.Pages.Paper
{
    public class ShareModel : PageModel
    {
        private readonly ApplicationDbContext _context;
        private readonly IAuthorizationService _authorizationService;
        private readonly PDFBuilder _pdfBuilder;
        private readonly TokenGenerator _tokenGenerator;

        public ShareModel(IAuthorizationService authorizationService, ApplicationDbContext context,
            PDFBuilder pdfBuilder, TokenGenerator tokenGenerator)
        {
            _authorizationService = authorizationService;
            _context = context;
            _pdfBuilder = pdfBuilder;
            _tokenGenerator = tokenGenerator;
        }

        [BindProperty] public Models.Paper Paper { get; set; }

        [TempData] public string StatusMessage { get; set; }
        
        [BindProperty] public string DownloadToken { get; set; }


        public async Task<IActionResult> OnGetAsync(int? id)
        {
            if (id == null)
            {
                return NotFound();
            }

            Paper = await _context.Papers.Include(paper => paper.Author).FirstOrDefaultAsync(m => m.ID == id);

            if (Paper == null)
            {
                return NotFound();
            }

            var authorizationResult = await _authorizationService
                .AuthorizeAsync(User, Paper, PaperOperations.Share);

            if (authorizationResult.Succeeded)
            {
                try
                {
                    var token = _tokenGenerator.GenerateTokenAsync(Paper);
                    DownloadToken = token;
                }
                catch (FileNotFoundException)
                {
                    return NotFound();
                }
            }
            else
            {
                StatusMessage = "Error: You are not authorized to view this file.";
            }

            return Page();
            
        }
    }
}