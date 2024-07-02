document.addEventListener('DOMContentLoaded', function () {
  const searchInput = document.getElementById('search');
  const clearSearchButton = document.getElementById('clear-search');
  const results = document.getElementById('results');
  const noResults = document.getElementById('no-results');
  const searchingAnimation = document.getElementById('searching-animation');
  const searchingText = document.getElementById('search-msg');

  let bookings = [];

  function fetchBookings() {
    fetch('/api/bookings')
      .then((response) => response.json())
      .then((data) => {
        bookings = data;
        performSearch();
      })
      .catch((error) => {
        console.error('Error fetching bookings:', error);
      });
  }

  function displayBookings(bookingsToDisplay) {
    results.innerHTML = '';
    if (bookingsToDisplay.length === 0) {
      noResults.classList.remove('hidden');
      searchingAnimation.classList.remove('hidden');
    } else {
      noResults.classList.add('hidden');
      searchingAnimation.classList.add('hidden');
      bookingsToDisplay.forEach((booking) => {
        const card = document.createElement('div');
        card.classList.add('booking-card');
        card.setAttribute('data-contact', booking.name);
        card.setAttribute('data-uuid', booking.uuid);

        const truncEmail = `${booking.email.split('@')[0].slice(0, 3)}...@${
          booking.email.split('@')[1]
        }`;
        const formattedDateOfBooking = new Date(
          booking.start_time
        ).toLocaleString();

        card.innerHTML = `
          <div class="card-content">
            <h3>${booking.room_name}</h3>
            <p>Start Time: <strong><u>${formattedDateOfBooking}</u></strong></p>
            <p>Booked By: ${booking.name}</p>
            </div>
            <div class="card-expanded hidden">
            <p>Booking's Email: ${truncEmail}</p>
            <p>Booking's Phone: ${booking.phone}</p>
            <p>Group Size: ${booking.group_size}</p>
            <hr>
            <p id="confirm-text">Is this the correct booking?</p>
            <button class="confirm-button">Sign Waiver(s)</button>
          </div>
        `;

        card.addEventListener('click', function (e) {
          if (!e.target.classList.contains('confirm-button')) {
            const expandedSection = this.querySelector('.card-expanded');
            const isExpanding = expandedSection.classList.contains('hidden');

            // Collapse all cards and remove reduced opacity
            document.querySelectorAll('.booking-card').forEach((otherCard) => {
              otherCard.querySelector('.card-expanded').classList.add('hidden');
              otherCard.classList.remove('reduced-opacity');
            });

            if (isExpanding) {
              // Expand this card and reduce opacity of others
              expandedSection.classList.remove('hidden');
              document
                .querySelectorAll('.booking-card')
                .forEach((otherCard) => {
                  if (otherCard !== this) {
                    otherCard.classList.add('reduced-opacity');
                  }
                });
            }
          }
        });

        card
          .querySelector('.confirm-button')
          .addEventListener('click', function (e) {
            e.stopPropagation();
            const companyName = 'lockedmanhattan';
            const uuid = card.getAttribute('data-uuid');
            const url = `https://fareharbor.com/waivers?shortname=${companyName}&bookingUuid=${uuid}/`;
            window.location.href = url;
          });

        results.appendChild(card);
      });
    }
    toggleClearButton();
  }

  function isSearchTermValid(term, value) {
    return (
      term.length >= Math.ceil(value.split(' ')[0].length / 2) &&
      value.toLowerCase().includes(term)
    );
  }

  function performSearch() {
    const searchValue = searchInput.value.trim();

    if (searchValue === '@all') {
      displayBookings(bookings);
      return;
    }

    const searchTerms = searchValue
      .toLowerCase()
      .split(' ')
      .filter((term) => term.length > 0);

    if (searchTerms.length === 0) {
      results.innerHTML = '';
      noResults.classList.add('hidden');
      searchingAnimation.classList.add('hidden');
      toggleClearButton();
      return;
    }

    searchingText.innerText =
      searchValue.length > 8
        ? 'This booking may not be under your current search text.'
        : 'Keep typing to find your booking...';

    const filteredBookings = bookings.filter((booking) =>
      searchTerms.every(
        (term) =>
          isSearchTermValid(term, booking.name) ||
          isSearchTermValid(term, booking.email) ||
          isSearchTermValid(term, booking.phone)
      )
    );

    displayBookings(filteredBookings);
  }

  function clearSearch() {
    searchInput.value = '';
    results.innerHTML = '';
    noResults.classList.add('hidden');
    searchingAnimation.classList.add('hidden');
    toggleClearButton();
  }

  function toggleClearButton() {
    clearSearchButton.classList.toggle(
      'hidden',
      searchInput.value.length === 0 && results.children.length === 0
    );
  }

  searchInput.addEventListener('input', performSearch);

  clearSearchButton.addEventListener('click', clearSearch);

  window.addEventListener('popstate', function (event) {
    if (event.state && event.state.search !== undefined) {
      searchInput.value = event.state.search;
      performSearch();
    }
  });

  const urlParams = new URLSearchParams(window.location.search);
  const initialSearch = urlParams.get('search');
  if (initialSearch) {
    searchInput.value = initialSearch;
  }

  fetchBookings();
});
