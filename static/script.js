gsap.registerPlugin(ScrollTrigger);
document.addEventListener('DOMContentLoaded', () => {
  const form = document.getElementById('contact-form');
  const contactTitle = document.getElementById('contact-title');
  const contactBtn = document.getElementById('contact-btn');
  const header = document.querySelector('header');
  const hero = document.querySelector('.hero');

  const parallaxElements = document.querySelectorAll('.parallax-element');

  document.querySelectorAll('a[href^="#"]').forEach((anchor) => {
    anchor.addEventListener('click', function (e) {
      e.preventDefault();
      document.querySelector(this.getAttribute('href')).scrollIntoView({
        behavior: 'smooth',
      });
    });
  });

  gsap.utils.toArray('section').forEach((section) => {
    if (section.id === 'contact') {
      gsap.from(section, {
        opacity: 0,
        y: 50,
        duration: 1,
        scrollTrigger: {
          trigger: section,
          start: 'top 80%',
          end: 'top 70%',
          scrub: 1,
        },
      });
      return;
    }
    gsap.from(section, {
      opacity: 0,
      y: 50,
      duration: 1,
      scrollTrigger: {
        trigger: section,
        start: 'top 80%',
        end: 'top 50%',
        scrub: 1,
      },
    });
  });

  gsap.utils.toArray('.project-card').forEach((card) => {
    gsap.from(card, {
      opacity: 0,
      y: 30,
      duration: 0.5,
      scrollTrigger: {
        trigger: card,
        start: 'top 90%',
        end: 'top 60%',
        scrub: 1,
      },
    });
  });

  async function post(path, data, method = 'post') {
    let { name, email, message } = data;
    try {
      const res = await fetch(path, {
        method: method,
        body: JSON.stringify({ name, email, message }),
        headers: {
          Accept: 'application/json',
          'Content-Type': 'application/json',
        },
      });
    } catch (err) {
      console.log(err.message);
    }
  }

  form.addEventListener('submit', (e) => {
    e.preventDefault();
    const formData = new FormData(form);
    const data = Object.fromEntries(formData);

    post('/emailmsg/', data);

    form.reset();

    contactTitle.innerText = 'Thank you for contacting me!';
    contactBtn.innerText = 'Send Another Message';
  });

  const heroObserver = new IntersectionObserver(
    ([entry]) => {
      if (!entry.isIntersecting) {
        header.classList.add('scrolled');
      } else {
        header.classList.remove('scrolled');
      }
    },
    { threshold: 0.1 }
  );

  heroObserver.observe(hero);

  hero.addEventListener('mousemove', (e) => {
    const { clientX, clientY } = e;
    const { offsetWidth, offsetHeight } = hero;

    parallaxElements.forEach((el) => {
      const depth = el.getAttribute('data-depth');
      const moveX = (clientX - offsetWidth / 2) * depth;
      const moveY = (clientY - offsetHeight / 2) * depth;

      setTimeout(() => {
        el.style.transform = `rotate(${moveX + (moveY % 30)}deg)`;
      }, 200);
    });
  });

  hero.addEventListener('mouseleave', () => {
    parallaxElements.forEach((el) => {
      el.style.transform = 'translate(0, 0)';
    });
  });
});
