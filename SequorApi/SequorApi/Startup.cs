using AutoMapper;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.HttpsPolicy;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using Microsoft.OpenApi.Models;
using SequorApi.Negocio;
using SequorApi.Servico;
using SequorApi.Suporte;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

namespace SequorApi
{
    public class Startup
    {
        readonly string LiberacaoOrigens = "_liberacaoOrigens";

        public Startup(IConfiguration configuration)
        {
            Configuration = configuration;
        }

        public IConfiguration Configuration { get; }

        // This method gets called by the runtime. Use this method to add services to the container.
        public void ConfigureServices(IServiceCollection services)
        {
            var mapperConfig = new MapperConfiguration(mc =>
            {
                mc.AddProfile(new MappingProfile());
            });

            services.AddCors(options =>
            {
                options.AddPolicy(name: LiberacaoOrigens,
                                  policy =>
                                  {
                                      policy.WithOrigins("http://localhost:3000").AllowAnyHeader().AllowAnyMethod().AllowAnyOrigin();
                                                          
                                  });
            });


            IMapper mapper = mapperConfig.CreateMapper();
            services.AddSingleton(mapper);
            

            services.AddControllers();
            services.AddSwaggerGen(c =>
            {
                c.SwaggerDoc("v1", new OpenApiInfo { Title = "SequorApi", Version = "v1" });
            });

            services.AddScoped<IGitHubNeg, GitHubNeg>();
            services.AddScoped<IGitHubService, GitHubService>();
            
            services.AddHttpClient();
        }

        // This method gets called by the runtime. Use this method to configure the HTTP request pipeline.
        public void Configure(IApplicationBuilder app, IWebHostEnvironment env)
        {
            if (env.IsDevelopment())
            {
                app.UseDeveloperExceptionPage();
                app.UseSwagger();
                app.UseSwaggerUI(c => c.SwaggerEndpoint("/swagger/v1/swagger.json", "SequorApi v1"));
            }

            app.UseHttpsRedirection();

            app.UseRouting();

            app.UseCors(LiberacaoOrigens);

            app.UseAuthorization();

            app.UseEndpoints(endpoints =>
            {
                endpoints.MapControllers();
            });
        }
    }
}


    
