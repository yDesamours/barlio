function previewPhoto(event) {
const boxSize = 400;
    const input = event.target;
    const preview = document.getElementById("photo_preview");
    const cropImage = document.getElementById("photo_crop");
    const myModal = new bootstrap.Modal(document.getElementById("modal"))
    let cropper = null;
  
    if (input.files && input.files[0]) {
      const  reader = new FileReader();
    
      myModal.show()

      reader.onload = function(e) {
        preview.src = e.target.result;
        preview.style.display = "block";

        cropImage.src = e.target.result;
        cropImage.style.display = "block";

        if (cropper) {
          cropper.destroy();
      }
      
      cropper = new Cropper(cropImage, {
        aspectRatio: 1,
        viewMode: 1,    
        maxCropBoxHeight : boxSize,
        maxCropBoxWidth : boxSize,
        minCropBoxWidth : boxSize,
        maxCropBoxHeight:boxSize,
        crop: function(event) {
          var croppedData = {
            x: event.detail.x,
            y: event.detail.y,
            width: event.detail.width,
            height: event.detail.height
          };
          console.log(croppedData);
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