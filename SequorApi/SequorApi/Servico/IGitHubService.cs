using SequorApi.Modelo.DTO;
using SequorApi.Modelo.GitHub;
using System.Collections.Generic;
using System.Threading.Tasks;

namespace SequorApi.Servico
{
    public interface IGitHubService
    {
        Task<PerfilGirHub> BuscaPerfilGitHubPorLogin(string login);

        Task<List<RepositorioGitHub>> BuscaListaRepositorio(string url);
    }
}
