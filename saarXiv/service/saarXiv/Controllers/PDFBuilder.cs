using System.Diagnostics;
using System.IO;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using saarXiv.Data;
using saarXiv.Models;

namespace saarXiv.Controllers
{
    public class PDFBuilder
    {
        private static string template = @"\documentclass{{article}}
\author{{{0}, {1}}}
\title{{{2}}}
\begin{{document}}
\maketitle

{3}
\end{{document}}
";

        private static string basedir = "./data";
        private static string pdflatex = "pdflatex_wrapper";
        private static string args = "-interaction nonstopmode -output-directory {0} {1}";
        private readonly ApplicationDbContext _context;

        public PDFBuilder(ApplicationDbContext context)
        {
            _context = context;
        }

        public async Task<ActionResult> GetPdfAsync(Paper paper)
        {
            return await Task.Run<ActionResult>(() => GetPdf(paper));
        }

        private ActionResult GetPdf(Paper paper)
        {
            string outdir = Path.Combine(basedir, paper.Author.UserName);
            string pdffile = $"{paper.SaarXivID}.pdf";
            string pdfpath = Path.Combine(outdir, pdffile);
            if (paper.NeedsCompilation)
            {
                if (File.Exists(pdfpath))
                {
                    File.Delete(pdfpath);
                }

                string texpath = Path.Combine(Path.GetTempPath(), $"{paper.SaarXivID}.tex");

                Directory.CreateDirectory(outdir);
                File.WriteAllText(texpath,
                    string.Format(template, paper.Author.Lastname, paper.Author.Firstname, paper.Title, paper.Content));

                var process = new Process
                {
                    StartInfo = {FileName = pdflatex, Arguments = string.Format(args, outdir, texpath)}
                };
                process.Start();
                process.WaitForExit();

                File.Delete(texpath);

                paper.NeedsCompilation = false;
                _context.Attach(paper).State = EntityState.Modified;
                _context.SaveChanges();
            }

            if (File.Exists(pdfpath))
            {
                var result = new FileStreamResult(File.OpenRead(pdfpath), "application/octet-stream")
                {
                    FileDownloadName = pdffile
                };
                return result;
            }

            throw new FileNotFoundException();
        }
    }
}
