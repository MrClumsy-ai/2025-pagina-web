document.addEventListener("DOMContentLoaded", () => {
  const verDetalles = document.querySelectorAll(".ver-detalles");
  const ocultarDetalles = document.querySelectorAll(".close-btn");
  verDetalles.forEach((button) => {
    button.addEventListener("click", () => {
      const targetDiv = button.closest("article").nextElementSibling;
      if (targetDiv.classList.contains("hidden")) {
        targetDiv.classList.replace("hidden", "flex");
      }
    });
  });
  ocultarDetalles.forEach((button) => {
    button.addEventListener("click", () => {
      button.closest(".modal-mascota").classList.replace("flex", "hidden");
    });
  });
});
