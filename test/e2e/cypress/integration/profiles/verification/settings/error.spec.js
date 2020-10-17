import {
  APP_URL,
  assertVerifiableAddress,
  gen,
  parseHtml,
  verifyHrefPattern,
} from '../../../../helpers'

context('Settings', () => {
  describe('error flow', () => {
    let identity
    before(() => {
      cy.deleteMail()
    })

    beforeEach(() => {
      identity = gen.identity()
      cy.register(identity)
      cy.deleteMail({ atLeast: 1 }) // clean up registration email

      cy.login(identity)
      cy.visit(APP_URL + '/settings')
    })

    it('is unable to verify the email address if the code is no longer valid', () => {
      const email = `not-${identity.email}`
      cy.get('#user-profile input[name="traits.email"]').clear().type(email)
      cy.get('#user-profile button[type="submit"]').click()
      cy.verifyEmailButExpired({ expect: { email } })
    })

    it('is unable to verify the email address if the code is incorrect', () => {
      const email = `not-${identity.email}`
      cy.get('#user-profile input[name="traits.email"]').clear().type(email)
      cy.get('#user-profile button[type="submit"]').click()

      cy.getMail().then((mail) => {
        const link = parseHtml(mail.body).querySelector('a')

        expect(verifyHrefPattern.test(link.href)).to.be.true

        cy.visit(link.href + '-not') // add random stuff to the confirm challenge
        cy.log(link.href)
        cy.session().then(assertVerifiableAddress({ isVerified: false, email }))
      })
    })

    xit('should not update the traits until the email has been verified and the old email has accepted the change', () => {
      // FIXME https://github.com/zzpu/openuser/issues/292
    })
  })
})
