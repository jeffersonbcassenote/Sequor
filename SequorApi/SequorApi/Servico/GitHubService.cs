using Microsoft.Extensions.Configuration;
using Microsoft.Net.Http.Headers;
using SequorApi.Modelo.DTO;
using SequorApi.Modelo.GitHub;
using System;
using System.Collections.Generic;
using System.Net.Http;
using System.Text.Json;
using System.Threading.Tasks;

namespace SequorApi.Servico
{
    public class GitHubService : IGitHubService
    {
        private readonly IConfiguration _configuration;
        private readonly IHttpClientFactory _httpClientFactory;
        private string _urlIntegracaoGitHub = null;

        private string URL_INTEGRACAO_GIT
        {
            get
            {
                return this._urlIntegracaoGitHub ?? (this._urlIntegracaoGitHub = _configuration.GetValue<string>("IntegracaoGitHub:URL"));
            }
        }

        public GitHubService(IConfiguration configuration, IHttpClientFactory httpClientFactory)
        {
            _configuration = configuration;
            _httpClientFactory = httpClientFactory;
        }

        public Task<List<RepositorioGitHub>> BuscaListaRepositorio(string url)
        {
            return RequisicaoAsync<List<RepositorioGitHub>>(url);
        }

        public Task<PerfilGirHub> BuscaPerfilGitHubPorLogin(string login)
        {
            string urlFormatada = URL_INTEGRACAO_GIT + login;

            return RequisicaoAsync<PerfilGirHub>(urlFormatada);
        }

        private async Task<T> RequisicaoAsync<T>(string url)
        {
            T retorno = default(T);

            var httpRequestMessage = new HttpRequestMessage(HttpMethod.Get, url)
            {
                Headers = {
                    { HeaderNames.Accept, "application/vnd.github.v3+json" },
                    { HeaderNames.UserAgent, "HttpRequests" }
                }
            };

            var httpClient = _httpClientFactory.CreateClient();
            var httpResponseMessage = await httpClient.SendAsync(httpRequestMessage);

            if (httpResponseMessage.IsSuccessStatusCode)
            {
                using var contentStream =
                    await httpResponseMessage.Content.ReadAsStreamAsync();

                retorno = await JsonSerializer.DeserializeAsync<T>(contentStream);
            }
            

            return retorno;
        }
    }
}
