let cropper = null;

function cropPhoto(event) {
const boxSize = 400;
    const input = event.target;
    const cropImage = document.getElementById("crop");
    const myModal = new bootstrap.Modal(document.getElementById("modal"))
  
    if (input.files && input.files[0]) {
      const  reader = new FileReader();
    
      myModal.show()

      reader.onload = function(e) {
        cropImage.src = e.target.result;
        cropImage.style.display = "block";

        if (cropper) {
          cropper.destroy();
      }
      
      cropper = new Cropper(cropImage, {
        aspectRatio: 1,
        viewMode: 1,    
        background:false,
        modal:false,
        maxCropBoxHeight : boxSize,
        maxCropBoxWidth : boxSize,
        minCropBoxWidth : boxSize,
        maxCropBoxHeight:boxSize,
        crop: function(event) {
          var croppedData = {
            x: event.detail.x,
            y: event.detail.y,
            width: boxSize,
            height: boxSize
          };
        }
      });
    }
    
    reader.readAsDataURL(input.files[0]);
  } else {
    preview.src = "#";
    preview.style.display = "none";
    cropImage.src = "#";
    cropImage.style.display = "none";
  }
}

function previewPhoto(){
  const preview = document.getElementById("preview");
  const croppedCanvas = cropper.getCroppedCanvas();
  preview.style.backgroundImage = 'url('+croppedCanvas.toDataURL('image/jpeg')+')'
  console.log(croppedCanvas)
}
