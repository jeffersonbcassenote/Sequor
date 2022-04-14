using System.Collections.Generic;

namespace SequorApi.Modelo.DTO
{
    public class PerfilGitDTO
    {
        public PerfilGitDTO()
        {
            ListaRepositorio = new List<RepositorioGitDTO>();
        }

        public string Nome { get; set; }
        public string Empresa { get; set; }
        public string Blog { get; set; }
        public string Login { get; set; }
        public string Localizacao{ get; set; }
        public string Biografia { get; set; }
        public string Email { get; set; }
        public string Avatar { get; set; }
        public string Twitter { get; set; }
        
        public List<RepositorioGitDTO> ListaRepositorio { get; set; }
    }
}
