using AutoMapper;
using SequorApi.Modelo.DTO;
using SequorApi.Modelo.GitHub;

namespace SequorApi.Suporte
{
    public class MappingProfile : Profile
    {
        public MappingProfile()
        {
            CreateMap<PerfilGirHub, PerfilGitDTO>()
                .ForMember(dest => dest.Avatar, map => map.MapFrom(scr => scr.avatar_url))
                .ForMember(dest => dest.Nome, map => map.MapFrom(scr => scr.name))
                .ForMember(dest => dest.Empresa, map => map.MapFrom(scr => scr.company))
                .ForMember(dest => dest.Blog, map => map.MapFrom(scr => scr.blog))
                .ForMember(dest => dest.Login, map => map.MapFrom(scr => scr.login))
                .ForMember(dest => dest.Localizacao, map => map.MapFrom(scr => scr.location))
                .ForMember(dest => dest.Biografia, map => map.MapFrom(scr => scr.bio))
                .ForMember(dest => dest.Email, map => map.MapFrom(scr => scr.email))
                .ForMember(dest => dest.Twitter, map => map.MapFrom(scr => scr.twitter_username)).ReverseMap();


            CreateMap<RepositorioGitHub, RepositorioGitDTO>()
                .ForMember(dest => dest.Url, map => map.MapFrom(scr => scr.html_url))
                .ForMember(dest => dest.Nome, map => map.MapFrom(scr => scr.name))
                .ForMember(dest => dest.Estrela, map => map.MapFrom(scr => scr.stargazers_count))
                .ForMember(dest => dest.Descricao, map => map.MapFrom(scr => scr.description)).ReverseMap();
        }
    }
}
