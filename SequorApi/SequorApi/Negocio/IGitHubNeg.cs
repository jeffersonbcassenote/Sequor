using SequorApi.Modelo.DTO;
using System.Threading.Tasks;

namespace SequorApi.Negocio
{
    public interface IGitHubNeg
    {
        Task<PerfilGitDTO> BuscaPerfilGitHubPorLogin(string login);
    }
}
