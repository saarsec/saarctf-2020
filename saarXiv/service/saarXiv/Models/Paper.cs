using System.ComponentModel.DataAnnotations;

namespace saarXiv.Models
{
    public class Paper
    {
        [Key] public int ID { get; set; }

        public virtual Author Author { get; set; }

        public string Title { get; set; }

        public string Content { get; set; }

        public bool UnderSubmission { get; set; }

        public bool NeedsCompilation { get; set; }

        public string SaarXivID
        {
            get { return string.Format("saarXiv_{0:00000}", ID); }
        }
    }
}