using AutoMapper;
using SequorApi.Modelo.DTO;
using SequorApi.Servico;
using System.Collections.Generic;
using System.Threading.Tasks;
using System.Linq;

namespace SequorApi.Negocio
{
    public class GitHubNeg : IGitHubNeg
    {
        private readonly IMapper _mapper;
        private readonly IGitHubService _gitHubService;

        public GitHubNeg(IGitHubService gitHubService, IMapper mapper)
        {
            _gitHubService = gitHubService;
            _mapper = mapper;
        }

        public async Task<PerfilGitDTO> BuscaPerfilGitHubPorLogin(string login)
        {
            var perfil = await _gitHubService.BuscaPerfilGitHubPorLogin(login);

            if (perfil is null)
                return null;

            var retorno = _mapper.Map<PerfilGitDTO>(perfil);

            if (!string.IsNullOrWhiteSpace(perfil.repos_url))
            {
                var lst = await _gitHubService.BuscaListaRepositorio(perfil.repos_url);

                lst?.OrderByDescending(o=> o.stargazers_count).ToList().ForEach(x => retorno.ListaRepositorio.Add(_mapper.Map<RepositorioGitDTO>(x)));
            }
            
            return retorno;

        }
    }
}
