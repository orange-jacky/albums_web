window.onload = function() {
  changeColor('content');
}

function changeColor(id) {
  var arrayli = document.getElementById(id)
    .getElementsByTagName('tr');
  var bool = true; //奇数行为true 
  var oldStyle; //保存原有样式 
  for (var i = 1; i < arrayli.length; i++) {
    //各行变色 
    if (bool === true) {
      arrayli[i].className = "change";
      bool = false;
    } else {
      arrayli[i].className = "";
      bool = true;
    }
    //划过变色 
    arrayli[i].onmouseover = function() {
      oldStyle = this.className;
      this.className = "change2"
    }
    arrayli[i].onmouseout = function() {
      this.className = oldStyle;
    }
  }
}