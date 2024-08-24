import {useEffect} from 'react';
import {useState} from 'react';
import './App.css'

const CLIENT_ID = "Ov23liOrQ94iafbT1wjK"
const URL = "localhost:8000"
const GIT_LOGIN_URL = "https://github.com/login/oauth"


function App(){
  const [isLogged, setIsLogged] = useState(false);
  const [number, setNumber] = useState<string>("");
  const [credentials, setCredentials] = useState(["", ""]);

  useEffect(()=> {
    const queryString = window.location.search;
    const urlParams = new URLSearchParams(queryString);
    
    var code = urlParams.get("code")
    if (code == null && localStorage.getItem("accessToken") === null){
      setIsLogged(false)
    }else{
      setIsLogged(true)
      if (localStorage.getItem("accessToken") === null){
        getAccessToken(code)
      } else {
        getUserData(localStorage.getItem("accessToken"))
      }
      const socket = new WebSocket("ws://" + URL + "/ws")
    
      socket.addEventListener("message", event => {
        setNumber(event.data);
      });
    }
  }, []);

  function loginWithGitHub() {
    window.location.assign(GIT_LOGIN_URL + "/authorize?client_id=" + CLIENT_ID + "&prompt=consent")
  }

  function logoutFromGitHub(){
    localStorage.removeItem("accessToken");
    window.location.assign("http://localhost:5173");
  }

  async function getAccessToken(code: string | null){
    await fetch("http://" + URL + "/getAccessToken", {
      method:"POST",
      body:JSON.stringify({
        id:CLIENT_ID, 
        url:GIT_LOGIN_URL + "/access_token",
        code:code
      })
    }).then((response) => {
          return response.json();
        }).then((data) => {if (localStorage.getItem("accessToken") === null){
          localStorage.setItem("accessToken", data.token)
          getUserData(data.token)
        }})
  }

  async function getUserData(token: string | null){
    await fetch("http://" + URL + "/getUserData", {
      method:"POST",
      body:JSON.stringify({
        token:token
      })
    }).then((response) => {
          return response.json();
        }).then((data) => {
          setCredentials([data.name, data.login])
        })
  }

  return (
    <div className='App'>
      <header className='App-header'>
        {!isLogged && <div className='container'>
          <div className='label'>
            Добро пожаловать в генератор случайных чисел!
          </div>
          <div className='logo'></div>
          <button className='button' onClick={loginWithGitHub}>
            Войти с GitHub
          </button>
        </div>}
        {isLogged && <div className='container'>
          <div className='label_greeting'>
            Добро пожаловать, {credentials[0]} или {credentials[1]}! Ваше счастливое число на эти 5 секунд!
          </div>
          <div className="circle" >
            {number != "" ? number : null}
          </div>
          <button className='button' onClick={logoutFromGitHub}>
            Выйти
          </button>
        </div>}
      </header>
    </div>
  );

}

export default App