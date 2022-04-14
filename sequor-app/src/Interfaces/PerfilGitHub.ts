
 import {RepositorioGitHub } from './RepositorioGitHub'

 

export interface PerfilGitHub {
    nome:string ;
    empresa: string ;
    blog: string ;
    login:string ;
    localizacao:string ;
    biografia: string ;
    email:string ;
    twitter: string ;
    avatar:string ;
    listaRepositorio:   Array<RepositorioGitHub>;
}





// export class PerfilGitHub {
//     _nome:string ;
//     empresa: string ;
//     blog: string ;
//     login:string ;
//     localizacao:string ;
//     biografia: string ;
//     email:string ;
//     twitter: string ;
//     _avatar:string ;
//     listaRepositorio:   Array<RepositorioGitHub>;

//     constructor(nome:string ,    empresa: string ,    blog: string ,    login:string ,    localizacao:string ,    biografia: string ,    email:string ,    twitter: string ,  avatar:string , listaRepositorio:   Array<RepositorioGitHub>) {
//       this._nome= nome;
//       this.empresa= empresa;
//       this.blog= blog;
//       this.login=login;
//       this.localizacao=localizacao;
//       this.biografia=biografia;
//       this.email= email;
//       this.twitter=twitter
//       this._avatar=avatar;
//       this.listaRepositorio= listaRepositorio;
//     }

//     get nome() {
//       return this._nome;
//     }
//     set nome(value) {
//       this._nome = value;
//     }

//     get avatar(){
//       return this._nome;
//   } 

// }


