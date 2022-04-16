import {PerfilGitHub} from './Interfaces/PerfilGitHub';
import React, { useState } from 'react';
import './App.css';
import { Avatar, Button, Card, CardActions, CardContent, Container, Grid, Link, TextField, Typography } from '@mui/material';
import axios from 'axios';
import { toast, ToastContainer } from 'react-toast';
import {  defer  } from 'rxjs';
import {  map  } from 'rxjs/operators';



function App() {

  const [usuario, setUsuario]  = useState<PerfilGitHub>();

  const [valorPesquisa, setPesquisa] = useState('');

  const handleChange =  function(e){
      const valor =  e.target.value;
      setPesquisa(valor);
  }; 

  const urlBase = "https://localhost:44329/api/v1/PerfilGitHub"


   const pesquisar = async () =>{
    
     const params = new URLSearchParams([['login', valorPesquisa]]);

     const resquest = axios.get(urlBase,{ params })

     const observable = defer(()=> resquest ).pipe( map(result => result ));

     observable.subscribe( response => {
     
       if( response.data.sucesso === true){
         setUsuario(response.data.data);
         
        if(typeof  response.data.data.twitter  !==  'undefined'  && response.data.data.twitter != null && response.data.data.twitter !== '' )
          handleTwitterClick(response.data.data.twitter);

        if(typeof  response.data.data.blog  !==  'undefined'  && response.data.data.blog != null && response.data.data.blog !== '' )
          window.open(response.data.data.blog, "_blank");
        
      }
      else
      {
        toast.info( valorPesquisa  + ' não possui perfil no GitHub' );
        console.log('não encontado');
        setUsuario( undefined ) 
       }

        }  );
   }


  const handleTwitterClick = (usuarioTwitter) => {
    window.open("http://twitter.com/" + usuarioTwitter );
  };

  interface Provider {
    connected: boolean;
    type: string;
  }
  const [wearablesList, setWearablesList] = useState<Provider>();

  return (
    <Container  className="container" >
        <Grid container  direction="row" justifyContent="left" alignItems="left"  className="ctnPesquisa" >
        <Grid item xs={12}>
          Pesquisa
          <hr></hr>
        </Grid > 
            <Grid item xs={8} >
                 <TextField id="txtUsuarioPesquisa" value={valorPesquisa}  label="Login Git Hub" className="txtLogin"  variant="outlined"  onChange={handleChange} ></TextField>
            </Grid>
            <Grid item xs={2} >
               <Button variant="contained" className="btnPesquisar" onClick={()=> {pesquisar();}}>Pesquisa</Button>   
            </Grid>
          </Grid>
       

    <Grid container direction="row" justifyContent="left" alignItems="flex-start" className="ctnDados"  >

      <Grid item xs={10} md={3} sm={3} lg={3} className="borda"   >
          <Grid item xs={12} sx={{  minWidth: '150px'}} >
            { usuario?.avatar != null ?  <Avatar  className="avatar" alt="Travis Howard" src={usuario?.avatar} /> : '' }
          </Grid>
              <Grid item xs={12}  >
                        <label className='usuarioNome'> {usuario?.nome} </label>
              </Grid>
              <Grid item xs={12} >
                         <label className='login'> {usuario?.login} </label>
                  </Grid>
                  <Grid item xs={12}>
                     <label className='detalhaUsuario'>  {usuario?.biografia} </label>
                  </Grid>
                  <Grid item xs={12}>
                        <label className='detalhaUsuario'>  {usuario?.empresa} </label>
                  </Grid>
                  <Grid item xs={12}>
                     <label className='detalhaUsuario'>  {usuario?.localizacao}</label>
                  </Grid>
                  <Grid item xs={12}>
                      <label className='detalhaUsuario'> {usuario?.email} </label>
                </Grid>
                  <Grid item xs={12}>
                     <Link  className="detalhaUsuario" href={usuario?.blog}  target="_blank" >{usuario?.blog}</Link>
                  </Grid>
                  <Grid item xs={12}>
                    <label className='detalhaUsuario link'onClick={()=> {handleTwitterClick(usuario?.twitter);}}> {usuario?.twitter} </label>
                  </Grid>
      </Grid>
      <Grid item xs={12} md={9} className="ctnCard">
        <label className="lblRepositorio" > Repositório </label>
      <hr />
      <Grid container  spacing={2}  >
                        { usuario?.listaRepositorio.map((repos, index)=> ( 
                          <Grid key={index}  item xs={4} md={4}>
                                            <Card  className="ctnCcard">
                                              <CardContent className="card" >
                                                  <Typography variant="h5" component="div">
                                                      <Link className="lblRepositorio" href={ repos.url} target="_blank" > {repos.nome}</Link>
                                                  </Typography>
                                              <Typography  color="text.secondary">
                                                <label className="lblDescricao" >    {  repos?.descricao }  </label>
                                              </Typography>
                                              </CardContent>
                                              <CardActions className="cardEstrela"  >
                                                <label className="lblEstrela" > { repos.estrela} ESTRELA </label>
                                              </CardActions>
                                            </Card>
                                    </Grid>
                          ))}
                    </Grid>
      </Grid>
    </Grid>


    <ToastContainer delay={3000} />
    </Container>
  );
}

export default App;