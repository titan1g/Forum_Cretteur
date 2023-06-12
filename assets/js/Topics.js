var topicval = document.getElementsByClassName("category");
var topic = document.getElementsByClassName("top");
var valueselect = document.getElementById("select");
var on = true;


function select() {
    for (let i = 0; i < topicval.length; i++) {
        str = topicval[i].value;
        var arr = str.split(",");
        on = true;
        arr.forEach(v => {
            if ((v == valueselect.value && on) || valueselect.value == "All the topics") {
                topic[i].classList.add("on");
                topic[i].classList.remove("off");
                on = false;
            } else if (on) {
                topic[i].classList.add("off");
                topic[i].classList.remove("on");
            }
        });
    }
}