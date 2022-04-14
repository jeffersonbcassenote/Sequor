using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using SequorApi.Modelo.DTO;
using SequorApi.Negocio;
using SequorApi.Servico;
using System;
using System.Threading.Tasks;

namespace SequorApi.Controllers
{
    [Route("api/v1/[controller]")]
    [ApiController]
    public class PerfilGitHubController : ControllerBase
    {
        private readonly IGitHubNeg _gitHubNeg;

        public PerfilGitHubController(IGitHubNeg gitHubNeg)
        {
            _gitHubNeg = gitHubNeg;
        }

        [HttpGet()]
        public async Task<ActionResult<PerfilGitDTO>> BuscaPerfilGitHub(string login)
        {
            if (string.IsNullOrWhiteSpace(login))
                return BadRequest();
            try
            {
                var retorno = await _gitHubNeg.BuscaPerfilGitHubPorLogin(login);

                var r = new
                {
                    sucesso = retorno != null,
                    data = retorno
                };

                return Ok(r);

            }
            catch
            {
                return BadRequest();
            }

        }
    }
}
