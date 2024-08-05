import {useEffect} from 'react';
import {useState} from 'react';
import './App.css'

const CLIENT_ID = "Ov23liOrQ94iafbT1wjK"

function App(){
  const [isLogged, setIsLogged] = useState(false);
  const [number, setNumber] = useState<{generated: string} | null>(null);
  
  useEffect(()=> {
    const queryString = window.location.search;
    const urlParams = new URLSearchParams(queryString);
    var code = urlParams.get("code")
    if (code == null && localStorage.getItem("accessCode") === null){
      setIsLogged(false)
    }else{
      setIsLogged(true)
      getAccessToken(code)
      if (localStorage.getItem("accessCode") === null){
        if (code != null){
          localStorage.setItem("accessCode", code)
        }
      }
      let interval = setInterval(async () => {
        await fetch("http://localhost:8000/getNumber", {
          method: "GET"
        })
        .then(res => res.json())
        .then(
          (result) => {
            setNumber(result);
          }
        );
      }, 500);
      return () => {
        clearInterval(interval);
      }
    }
  }, []);

  function loginWithGitHub() {
    window.location.assign("https://github.com/login/oauth/authorize?client_id=" + CLIENT_ID //+ "&prompt=consent"
      )
  }

  function logoutFromGitHub(){
    localStorage.removeItem("accessCode");
    window.location.assign("http://localhost:5173");
  }

  async function getAccessToken(code: string | null){
    console.log("lel")
    await fetch("http://localhost:8000/getAccessToken", {
      method:"POST",
      body:JSON.stringify({
        id:CLIENT_ID, 
        url:"https://github.com/login/oauth/access_token",
        code:code
      })
    }).then((response) => {
          return response.json();
        }).then((data) => {
           console.log(data)
        })
  }

  // async function getAccessToken(code){
  //   const params = "?client_id=" + CLIENT_ID + "&client_secret=" + CLIENT_SECRET + "&code=" + code;
  //   console.log(params);
  //   const link = "https://github.com/login/oauth/access_token" + params;
  //   await fetch(link, {
  //     mode: 'cors',
  //     method: "POST",
  //     headers: {
  //       "Accept": "application/json",
  //       "Access-Control-Allow-Origin": "*"
  //     }

  //   }).then((response) => {
  //     return response.json();
  //   }).then((data) => {
  //      console.log(data)
  //   })
  // }

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
          <div className='label'>
            Ваше счастливое число на эти 5 секунд!
          </div>
          <div className="circle" >
            {number != null ? number.generated : null}
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