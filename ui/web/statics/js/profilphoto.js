let cropper = null;
let croppedData

function openModal(){
  const myModal = new bootstrap.Modal(document.getElementById("modal"))
  myModal.show();
}

function closeModal(){
  const myModal = new bootstrap.Modal(document.getElementById("modal"))
  myModal.hide();
}

function editPhoto() {
  const boxSize = 400;
  const cropImage = document.getElementById("crop");
  
  openModal()

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
    croppedData = {
      x: event.detail.x,
      y: event.detail.y,
    };
  }
});
}

function changePhoto(event){
  const file = document.getElementById("profil_picture");
  const cropImage = document.getElementById("crop");
  const reader = new FileReader()

  if (!file.files[0]){
    return
  }

  reader.onload = function(event){
    cropImage.src = event.target.result;
    openModal()
  }

  reader.readAsDataURL(file.files[0])
}


function previewPhoto(){
  const preview = document.getElementById("preview");
  const cropImage = document.getElementById("crop");

  preview.src = cropImage.src
  preview.style.left=`-${croppedData.x}px`
  preview.style.top=`-${croppedData.y}px`
  console.log(croppedCanvas)
}
