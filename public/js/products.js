let cartIcon = document.querySelector('#cart-icon');
let cart = document.querySelector('.cart');
let closeCart = document.querySelector('#close-cart');

// open cart
cartIcon.onclick = () => {
    cart.classList.add('active');
}
// close cart
closeCart.onclick = () => {
    cart.classList.remove('active');
}

let com = document.querySelector('#com');
let comments = document.querySelector('.comments');
let closeCom = document.querySelector('#close-com');
// open comments
//com.onclick = () => {
//    comments.classList.add('active');
//}
// close comments
// closeCom.onclick = () => {
//     comments.classList.remove('active');
// }



function showColors(sizeName) {

    var listColor = sizeName.parentElement.parentElement.children[0].children;
    var sizeList = sizeName.parentElement.children
    for (const node of sizeList) {
        if (node.innerText == sizeName.innerText) {
            node.style.background = "black"
        } else {
            node.style.background = null
        }

    }
    for (const node of listColor) {
        if (node.id === sizeName.innerText) {
            node.style.display = "block";
        } else {
            node.style.display = "none";
            node.style.scale = 1
        }
    }

}
function checkColor(colorIn) {
    var listColor = colorIn.parentElement.children
    for (const colors of listColor) {
        if (colorIn == colors) {
            colors.style.scale = 1.3
        } else {
            colors.style.scale = 1
        }
    }
}