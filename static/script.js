document.addEventListener('DOMContentLoaded', function () {
  const searchInput = document.getElementById('search');
  const clearSearchButton = document.getElementById('clear-search');
  const results = document.getElementById('results');
  const noResults = document.getElementById('no-results');
  const searchingAnimation = document.getElementById('searching-animation');
  const seeAllBTN = document.getElementById('see-all');
  const searchingText = document.getElementById('search-msg');

  let contacts = [
    {
      contact: {
        email: 'yellow@jenny.com',
        is_subscribed_for_email_updates: false,
        language: 'en-us',
        name: 'Jenny Dude',
        normalized_phone: '+17857760168',
        phone: '+1-443-222-1100',
        phone_country: 'US',
      },
      uuid: '631f6877-319c-4602-bb5c-d09168e9c91e',
    },
    {
      contact: {
        email: 'blue@Bill.com',
        is_subscribed_for_email_updates: false,
        language: 'en-us',
        name: 'Bill Denning',
        normalized_phone: '+17853417421',
        phone: '+1-443-222-1100',
        phone_country: 'US',
      },
      uuid: '631f6877-319c-4602-bb5c-d09168e9c91e',
    },
    {
      contact: {
        email: 'Green@Alen.com',
        is_subscribed_for_email_updates: false,
        language: 'en-us',
        name: 'Allen Currs',
        normalized_phone: '+10232394290',
        phone: '+1-443-222-1100',
        phone_country: 'US',
      },
      uuid: '631f6877-319c-4602-bb5c-d09168e9c91e',
    },
    {
      contact: {
        email: 'orange@Jim.com',
        is_subscribed_for_email_updates: false,
        language: 'en-us',
        name: 'Allen Noors',
        normalized_phone: '+10232394290',
        phone: '+1-443-222-1100',
        phone_country: 'US',
      },
      uuid: '631f6877-319c-4602-bb5c-d09168e9c91e',
    },
  ];

  function renderResults(filteredContacts) {
    results.innerHTML = '';
    if (filteredContacts.length == 0) {
      noResults.classList.remove('hidden');
      searchingAnimation.classList.remove('hidden');
    } else {
      noResults.classList.add('hidden');
      searchingAnimation.classList.add('hidden');
      filteredContacts.forEach((booking) => {
        let {
          contact: { name, email, phone },
          uuid,
        } = booking;

        let [username, domain] = email.split('@');
        let truncEmail = `${username.slice(0, 3)}...@${domain}`;

        const li = document.createElement('li');
        li.textContent = `${name} (email: ${truncEmail})`;
        li.setAttribute('data-contact', name);
        li.setAttribute('data-uuid', uuid);
        results.appendChild(li);
      });
    }
  }

  function isSearchTermValid(term, value) {
    return (
      term.length >= Math.ceil(value.split(' ')[0].length / 2) &&
      value.toLowerCase().includes(term)
    );
  }

  function performSearch() {
    const searchTerms = searchInput.value
      .toLowerCase()
      .split(' ')
      .filter((term) => term.length > 0);

    if (searchTerms.length === 0) {
      results.innerHTML = '';
      noResults.classList.add('hidden');
      searchingAnimation.classList.add('hidden');
      return;
    }

    if (searchInput.value.length > 8) {
      searchingText.innerText =
        'This booking may not be under your current search text.';
    } else {
      searchingText.innerText = 'Keep typing to find your booking...';
    }

    const filteredContacts = contacts.filter((booking) =>
      searchTerms.every(
        (term) =>
          isSearchTermValid(term, booking.contact.name) ||
          isSearchTermValid(term, booking.contact.email) ||
          isSearchTermValid(term, booking.contact.normalized_phone)
      )
    );

    renderResults(filteredContacts);
  }

  function clearSearch() {
    searchInput.value = '';
    results.innerHTML = '';
    noResults.classList.add('hidden');
    searchingAnimation.classList.add('hidden');
    clearSearchButton.classList.add('hidden');
  }

  seeAllBTN.addEventListener('click', () => {
    clearSearchButton.classList.remove('hidden');
    seeAllBTN.classList.add('hidden');

    history.pushState({ search: searchInput.value }, '', `?search=`);

    renderResults(contacts);
  });

  searchInput.addEventListener('input', function () {
    // seeAllBTN.classList.add('hidden');
    performSearch();
    clearSearchButton.classList.toggle(
      'hidden',
      searchInput.value.length === 0
    );
    seeAllBTN.classList.toggle('hidden', searchInput.value.length > 0);

    history.pushState(
      { search: searchInput.value },
      '',
      `?search=${encodeURIComponent(searchInput.value)}`
    );
  });

  clearSearchButton.addEventListener('click', function () {
    seeAllBTN.classList.remove('hidden');
    clearSearch();
    history.pushState({ search: '' }, '', window.location.pathname);
  });

  results.addEventListener('click', function (event) {
    const target = event.target;

    searchInput.value = target.getAttribute('data-contact');

    if (target.tagName === 'LI' && !target.querySelector('.confirm-button')) {
      const uuid = target.getAttribute('data-uuid');

      // Remove existing confirm buttons
      results
        .querySelectorAll('.confirm-button')
        .forEach((button) => button.remove());

      const confirmButton = document.createElement('button');
      const companyName = 'lockedmanhattan';
      confirmButton.textContent = 'Yes!';
      confirmButton.classList.add('confirm-button');
      confirmButton.addEventListener('click', function (e) {
        e.stopPropagation();
        const url = `https://fareharbor.com/waivers?shortname=${companyName}&bookingUuid=${uuid}/`;
        window.location.href = url;
      });

      const confirmText = document.createElement('p');
      confirmText.textContent = 'Are you sure this is the right booking?';

      target.appendChild(confirmText);
      target.appendChild(confirmButton);
    }
  });

  window.addEventListener('popstate', function (event) {
    if (event.state && event.state.search !== undefined) {
      searchInput.value = event.state.search;
      performSearch();
      clearSearchButton.classList.toggle(
        'hidden',
        searchInput.value.length === 0
      );
    }
  });

  const urlParams = new URLSearchParams(window.location.search);
  const initialSearch = urlParams.get('search');
  if (initialSearch) {
    searchInput.value = initialSearch;
    performSearch();
    clearSearchButton.classList.remove('hidden');
    seeAllBTN.classList.add('hidden');
  }
});
