document.addEventListener('DOMContentLoaded', function () {
  const searchInput = document.getElementById('search');
  const clearSearchButton = document.getElementById('clear-search');
  const results = document.getElementById('results');
  const noResults = document.getElementById('no-results');
  const searchingAnimation = document.getElementById('searching-animation');
  const searchingText = document.getElementById('search-msg');
  const roomNameFilterBtns = document.querySelectorAll('.room-filter');
  let scrollEvent;

  roomNameFilterBtns.forEach((elem) => {
    elem.addEventListener('click', () => {
      searchInput.value = elem.innerText;
      clearFilterBtnSelected();
      elem.classList.add('selected');
      performSearch();
    });
  });

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

  function formattedEmail(email) {
    let emailArr = email.split('@');
    let [addr, domain] = emailArr;
    let truncChar = '.';
    let splitIdx = Math.round(addr.length / 2);
    let truncAddr = addr.slice(0, splitIdx);
    let truncation = truncChar.repeat(addr.length - splitIdx);

    return `${truncAddr}${truncation}@${domain}`;
  }

  function formattedDate(dateString) {
    return new Date(dateString).toLocaleString();
  }

  function formattedBookingName(bookingName) {
    let [fName, lName] = bookingName.split(' ');

    return `${fName} ${lName.slice(0, 1)}.`;
  }

  function displayBookings(bookingsToDisplay, displayOne = false) {
    results.innerHTML = '';
    if (bookingsToDisplay.length === 0) {
      noResults.classList.remove('hidden');
      searchingAnimation.classList.remove('hidden');
    } else {
      noResults.classList.add('hidden');
      searchingAnimation.classList.add('hidden');

      bookingsToDisplay.forEach((booking, idx) => {
        const card = document.createElement('div');
        card.classList.add('booking-card');
        card.setAttribute('data-contact', booking.name);
        card.setAttribute('data-uuid', booking.uuid);

        const truncEmail = formattedEmail(booking.email);
        const formattedDateOfBooking = formattedDate(booking.start_time);
        const formattedTruncBookingName = formattedBookingName(booking.name);
        card.innerHTML = `
          <div class="card-content">
            <h3>${idx + 1}: ${
          booking.room_name
        }: ${formattedTruncBookingName}</h3>
            <p>Start Time: <strong>${formattedDateOfBooking}</strong></p>

            </div>
            <div class="card-expanded ${displayOne ? 'hidden' : ''}">
            <p><strong>Booking's Email: </strong>${truncEmail}</p>
            <p>Group Size: ${booking.group_size}</p>
            <hr>
            <p id="confirm-text">Is this the correct booking?</p>
            <button class="confirm-button">Sign Waiver(s)</button>
          </div>
        `;

        card.addEventListener('click', function (e) {
          if (!e.target.classList.contains('reduced-opacity')) {
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
      term.length >= Math.ceil(value.split(' ')[0].length / 3) &&
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
      clearSearch();
      return;
    }

    searchingText.innerText =
      searchValue.length > 8
        ? 'Booking may be under a different name, email, or room name'
        : 'Keep typing to find your booking...';

    const filteredBookings = bookings.filter((booking) =>
      searchTerms.every(
        (term) =>
          isSearchTermValid(term, booking.name) ||
          isSearchTermValid(term, booking.email) ||
          isSearchTermValid(term, booking.room_name)
      )
    );

    displayBookings(filteredBookings, filteredBookings.length > 1);
  }

  function clearFilterBtnSelected() {
    roomNameFilterBtns.forEach((filterBtn) => {
      filterBtn.classList.remove('selected');
    });
  }

  function clearSearch() {
    searchInput.value = '';
    results.innerHTML = '';
    noResults.classList.add('hidden');
    searchingAnimation.classList.add('hidden');
    toggleClearButton();
    clearFilterBtnSelected();
    searchInput.focus = true;
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
