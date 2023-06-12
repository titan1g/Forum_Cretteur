var btn = document.querySelector(".info");
var toshow = document.querySelector(".hide");
// var navbar = document.getElementsByClassName("header")[0];

var test = true;

toshow.classList.add("off");
if (btn != null) {
    btn.addEventListener("mouseover", event => {
        test = false;
        toshow.classList.add("on");
        toshow.classList.remove("off");
    });
    btn.addEventListener("mouseout", event => {
        test = true;
        setTimeout(function() {
            if (test) {
                toshow.classList.add("off");
                toshow.classList.remove("on");
            }
        }, 2000);
    });
}

// window.onscroll = () => {
//     console.log("a")
//     if (window.scrollY > 1) {
//       navbar.classList.add("u-scrolled");
//       console.log("b")
//     } else {
//       navbar.classList.remove("u-scrolled");
//       console.log("c")
//     }
//   };
