document.addEventListener("DOMContentLoaded", () => {
  document.querySelectorAll(".ver-detalles").forEach((button) => {
    button.addEventListener("click", () => {
      const targetDiv = button.closest("article").nextElementSibling;
      if (targetDiv.classList.contains("hidden")) {
        targetDiv.classList.replace("hidden", "flex");
      }
    });
  });
  document.querySelectorAll(".close-btn").forEach((button) => {
    button.addEventListener("click", () => {
      button.closest(".modal-mascota").classList.replace("flex", "hidden");
    });
  });
});
