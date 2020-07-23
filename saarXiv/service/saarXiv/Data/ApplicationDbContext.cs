using Microsoft.AspNetCore.Identity.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore;
using saarXiv.Models;

namespace saarXiv.Data
{
    public class ApplicationDbContext : IdentityDbContext<Author>
    {
        public ApplicationDbContext(DbContextOptions<ApplicationDbContext> options)
            : base(options)
        {
            Database.EnsureCreated();
        }

        public DbSet<Paper> Papers { get; set; }
    }
}