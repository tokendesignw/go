document.addEventListener("DOMContentLoaded", function () {
    // Для каждого проекта добавляем событие на открытие
    const projectLinks = document.querySelectorAll('.project-link');
    projectLinks.forEach(function (link) {
        link.addEventListener('click', function (e) {
            e.preventDefault();
            const projectId = link.getAttribute('data-id');
            window.location.href = `/project/${projectId}`;
        });
    });
});
