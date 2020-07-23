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
    public class DownloadModel : PageModel
    {
        private readonly ApplicationDbContext _context;
        private readonly IAuthorizationService _authorizationService;
        private readonly PDFBuilder _pdfBuilder;
        private readonly TokenValidator _tokenValidator;

        public DownloadModel(IAuthorizationService authorizationService, ApplicationDbContext context,
            PDFBuilder pdfBuilder, TokenValidator tokenValidator)
        {
            _authorizationService = authorizationService;
            _context = context;
            _pdfBuilder = pdfBuilder;
            _tokenValidator = tokenValidator;
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
                .AuthorizeAsync(User, Paper, PaperOperations.Download);

            if (authorizationResult.Succeeded)
            {
                try
                {
                    var pdfFile = await _pdfBuilder.GetPdfAsync(Paper);
                    return pdfFile;
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

        public async Task<IActionResult> OnPostAsync(int? id)
        {
            
            if (string.IsNullOrEmpty(DownloadToken))
            {
                return Page();
            }
            
            if (id == null)
            {
                return NotFound();
            }

            Paper = await _context.Papers.Include(paper => paper.Author).FirstOrDefaultAsync(m => m.ID == id);

            if (Paper == null)
            {
                return NotFound();
            }

            var result = await _tokenValidator.ValidateAsync(Paper, DownloadToken);
            if (result.Succeeded)
            {
                try
                {
                    var pdfFile = await _pdfBuilder.GetPdfAsync(Paper);
                    return pdfFile;
                }
                catch (FileNotFoundException)
                {
                    return NotFound();
                }    
            }
            else
            {
                StatusMessage = $"Error: {result.Error}";
                return Page();
            }
            
        }
    }
}