const login = document.querySelector('.login');
const loginLink = document.querySelector('.login-link');
const registerLink = document.querySelector('.login-register');
const headerloginIcon = document.querySelector('.headerlogin-icon');
const closeLogin = document.querySelector('.close-login');

registerLink.addEventListener('click', ()=> {
    login.classList.add('active');
});

loginLink.addEventListener('click', ()=> {
    login.classList.remove('active');
});

headerloginIcon.addEventListener('click', ()=> {
    login.classList.add('active-popup');
});

closeLogin.addEventListener('click', ()=> {
    login.classList.remove('active-popup');
<<<<<<< HEAD
});


let pass = document.querySelector('#pass');
let password = document.querySelector('#password');
// open change password
pass.onclick = () => {
    password.classList.add('active');
}
=======
});
>>>>>>> beb469b2fdc40db0cf219a3937b9a57bf1959a5b
