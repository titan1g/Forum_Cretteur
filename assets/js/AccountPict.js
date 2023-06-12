function blobToDataURL(blob) {
  console.log(typeof blob)
    var fileReader = new FileReader();
    fileReader.onload = function(e) {callback(e.target.result);}
    fileReader.readAsDataURL(blob);
    console.log(blob)
  }