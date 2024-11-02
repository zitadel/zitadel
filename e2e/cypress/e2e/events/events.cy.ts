describe('events', () => {
  beforeEach(() => {
    cy.context().as('ctx');
  });

  it('events can be filtered', () => {
    const eventTypeEnglish = 'Instance added';
    cy.visit('/instance?id=events');
    cy.get('[data-e2e="event-type-cell"]').should('have.length', 20);
    cy.get('[data-e2e="open-filter-button"]').click();
    cy.get('[data-e2e="event-type-filter-checkbox"]').click();
    cy.get('mat-select[name="eventTypesList"]').click();
    cy.contains('mat-option', eventTypeEnglish).click();
    cy.get('body').type('{esc}');
    cy.contains('mat-select', 'Descending').click();
    cy.contains('mat-option', 'Descending').click();
    cy.get('[data-e2e="filter-finish-button"]').click();
    cy.get('[data-e2e="event-type-cell"]').should('have.length', 1);
  });
});
