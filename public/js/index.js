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


function logoutUser() {

    fetch("/user/logout", {
        method: "GET"
    }).then(re => re.ok).then(d => {
        window.location.reload()
    })

}


function loginIndex() {
    const username = document.getElementById("usernameLogin").value
    const password = document.getElementById("passwordLogin").value
    fetch("/user/login", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: `{"username":"${username}","password":"${password}"}`,
    }).then(Response => Response.json()).then(data => {

        if (data != "ok") {
            document.getElementById("spanErrorLogin").innerText = data
        } else {
                window.location.reload()
        }
    })

}

function ClearErrorSpan() {
    document.getElementById("spanErrorLogin").innerText = ""
    document.getElementById("spanErrorResister").innerText = ""
}



function PostSettings() {

    const firstname = document.getElementById("firstname").value + " " + document.getElementById("lastname").value
    console.log(firstname)
    fetch("/user/name", {
        method: "POST",
        headers: {
            "Content-Type": "application/x-www-form-urlencoded",
        },
        body: `firstName=${firstname}`,
    }).then(Response => Response.text()).then(data => {
        console.log(data)
    })


    const oldPassowrd = document.getElementById("currentpassword").value
    console.log(oldPassowrd)
    const newpassword = document.getElementById("newpassword").value
    console.log(newpassword)
    fetch("/user/password", {
        method: "POST",
        headers: {
            "Content-Type": "application/text",
        },
        body: `{"oldPassowrd":"${oldPassowrd}","newpassword":"${newpassword}"}`,
    }).then(Response => Response.json()).then(data => {
    })
}

function GetSetting() {

    fetch("/user/name", {
        method: "GET"
    }).then(Response => Response.json()).then(data => {
        if (data.length != 0) {
            let name = data.split(" ")
            document.getElementById("last-name").value = name[1]
            document.getElementById("first-name").value = name[0]
            console.log(name.length)
        }
    })
}
function postcomment() {

    const commentData = {
        comment: "dawsasd",
        stars: 2,
        container: "newContainer",
        kind: "short",
        idmodel: 1
    };

    fetch("/user/commint", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(commentData),
    })
        .then(response => response.json())
        .then(result => {
            // Handle the response result
            console.log(result);
        })
}