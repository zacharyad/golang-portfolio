document.addEventListener('DOMContentLoaded', function () {
  const searchInput = document.getElementById('search');
  const clearSearchButton = document.getElementById('clear-search');
  const results = document.getElementById('results');
  const noResults = document.getElementById('no-results');
  const searchingAnimation = document.getElementById('searching-animation');
  const seeAllBTN = document.getElementById('see-all');

  let contacts = [
    { name: 'Alice Johnson', uuid: '631f6877-319c-4602-bb5c-d09168e9c91e' },
    { name: 'Bob Smith', uuid: 'b2c3d4e5' },
    { name: 'Charlie Brown', uuid: 'c3d4e5f6' },
    { name: 'Diana Ross', uuid: 'd4e5f6g7' },
    { name: 'Ethan Hunt', uuid: 'e5f6g7h8' },
    { name: 'Fiona Apple', uuid: 'f6g7h8i9' },
    { name: 'George Michael', uuid: 'g7h8i9j0' },
    { name: 'Alice Montana', uuid: 'h8i9j0k1' },
    { name: 'Ian McKellen', uuid: 'i9j0k1l2' },
    { name: 'Julia Roberts', uuid: 'j0k1l2m3' },
    { name: 'Kevin Bacon', uuid: 'k1l2m3n4' },
    { name: 'Alice Dern', uuid: 'l2m3n4o5' },
    { name: 'Michael Jordan', uuid: 'm3n4o5p6' },
    { name: 'Nancy Drew', uuid: 'n4o5p6q7' },
    { name: 'Oscar Wilde', uuid: 'o5p6q7r8' },
    { name: 'Alice Cruz', uuid: 'p6q7r8s9' },
    { name: 'Quentin Tarantino', uuid: 'q7r8s9t0' },
    { name: 'Rachel Green', uuid: 'r8s9t0u1' },
    { name: 'Steve Jobs', uuid: 's9t0u1v2' },
    { name: 'Tina Turner', uuid: 't0u1v2w3' },
  ];

  function renderResults(filteredContacts) {
    results.innerHTML = '';
    if (filteredContacts.length == 0) {
      noResults.classList.remove('hidden');
      searchingAnimation.classList.remove('hidden');
    } else {
      noResults.classList.add('hidden');
      searchingAnimation.classList.add('hidden');
      filteredContacts.forEach((contact) => {
        const li = document.createElement('li');
        li.textContent = `${contact.name} (UUID: ${contact.uuid})`;
        li.setAttribute('data-contact', contact.name);
        li.setAttribute('data-uuid', contact.uuid);
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

    const filteredContacts = contacts.filter((contact) =>
      searchTerms.every(
        (term) =>
          isSearchTermValid(term, contact.name) ||
          isSearchTermValid(term, contact.uuid)
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
