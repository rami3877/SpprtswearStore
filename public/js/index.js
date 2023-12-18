
// Object.getOwnPropertyNames(data).forEach((val, idx, array) => {
//     console.log(`${val} -> ${data[val]}`);
// });
// console.log(data["men"])


fetch("/AllContainerAndKind", {
    method: "GET"
}).then(res => res.json()).then(data => {
    var listallContainer = Object.getOwnPropertyNames(data);
    listallContainer.forEach(element => {
        var li = document.createElement("li")
        var a = document.createElement("a")
        a.innerText = element
        a.setAttribute("href", "/" + element)
        li.appendChild(a)
        document.getElementById("navlist").appendChild(li)

        li = document.createElement("li")
        a = document.createElement("a")
        a.innerText = element
        a.setAttribute("href", "/" + element)
        li.appendChild(a)
        document.getElementById("main-contact").appendChild(li)
    });
})

function loginIndex() {
    const username = document.getElementById("usernameLogin").value
    const password = document.getElementById("passwordLogin").value
    console.log(username)
    console.log(password)
    fetch("/user/login", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: `{"username":"${username}","password":"${password}"}`,
    }).then(Response => Response.text()).then(data => {
        if (data != "ok") {
            document.getElementById("spanErrorLogin").innerText = data

        } else {
            document.location("/")
        }
    }
    )

}

function ClearErrorSpan() {
    document.getElementById("spanErrorLogin").innerText = ""
    document.getElementById("spanErrorResister").innerText = ""
}


function Resister() {
    const username = document.getElementById("usernameResister").value
    const email = document.getElementById("emailResister").value
    const password = document.getElementById("passwordResister").value
    console.log(username)
    console.log(email)
    console.log(password)
    fetch("/user/register", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: `{"username":"${username}","password":"${password}","email":"${email}"}`,
    }).then(Response => Response.text()).then(data => {
        if (data != "Create") {
            document.getElementById("spanErrorResister").innerText = data
        } else {
            document.location("/")
        }

    }
    )
}