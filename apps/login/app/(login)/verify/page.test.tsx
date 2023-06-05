import RootLayout from '#/app/layout';
import Page from './page'
import { RenderResult, render, screen, waitFor, renderHook } from '@testing-library/react';
import { act } from 'react-dom/test-utils';

describe('/login/verify', () => {

    const noUserIDError = "No userId provided!"
    const noCodeError = "No code provided!"
    const codeLabelText = "Code"

    describe.each`
        userId          | code            | directSubmit    | expectRenderedUserIDError | expectRenderedCodePrefilled               | expectRenderedCodeError
        | ${"123"}      | ${"xyz"}        | ${true}         | ${false}                  | ${"" /*TODO: We should expect "xyz"*/}    | ${false}
        | ${"123"}      | ${""}           | ${true}         | ${false}                  | ${""}                                     | ${true}
        | ${"123"}      | ${undefined}    | ${true}         | ${false}                  | ${""}                                     | ${true}
        | ${"123"}      | ${"xyz"}        | ${false}        | ${false}                  | ${"xyz"}                                  | ${false}
        | ${"123"}      | ${""}           | ${false}        | ${false}                  | ${""}                                     | ${false}
        | ${"123"}      | ${undefined}    | ${false}        | ${false}                  | ${""}                                     | ${false}
        | ${""}         | ${"xyz"}        | ${true}         | ${true}                   | ${false}                                  | ${false}
        | ${""}         | ${""}           | ${true}         | ${true}                   | ${false}                                  | ${false}
        | ${""}         | ${undefined}    | ${true}         | ${true}                   | ${false}                                  | ${false}
        | ${""}         | ${"xyz"}        | ${false}        | ${true}                   | ${false}                                  | ${false}
        | ${""}         | ${""}           | ${false}        | ${true}                   | ${false}                                  | ${false}
        | ${""}         | ${undefined}    | ${false}        | ${true}                   | ${false}                                  | ${false}
        | ${undefined}  | ${"xyz"}        | ${true}         | ${true}                   | ${false}                                  | ${false}
        | ${undefined}  | ${""}           | ${true}         | ${true}                   | ${false}                                  | ${false}
        | ${undefined}  | ${undefined}    | ${true}         | ${true}                   | ${false}                                  | ${false}
        | ${undefined}  | ${"xyz"}        | ${false}        | ${true}                   | ${false}                                  | ${false}
        | ${undefined}  | ${""}           | ${false}        | ${true}                   | ${false}                                  | ${false}
        | ${undefined}  | ${undefined}    | ${false}        | ${true}                   | ${false}                                  | ${false}
        `(`With code=$code, submit=$submit and userId=$userId`, ({ userId, code, directSubmit, expectRenderedUserIDError, expectRenderedCodePrefilled, expectRenderedCodeError }) => {

        let renderResult: RenderResult;
        beforeEach(async () => {
            await act(async () => {
                renderResult = render(await RootLayout({
                    children: await Page({
                        searchParams: {
                            code: code,
                            submit: directSubmit,
                            userID: userId,
                        }
                    })
                }))
                // TODO: Replace the above syntax for awaiting the JSX.Element once https://github.com/DefinitelyTyped/DefinitelyTyped/pull/65135 is released
                // renderResult = render(<Page searchParams={{
                //      code: code,
                //      submit: submit,
                //      userID: userId,
                //  }} />)
            })
        })

        it(`should have rendered`, () => {
            expect(renderResult.container.firstChild).toBeDefined();
        })
        describe(`With expectRenderedUserIDError=${expectRenderedUserIDError}`, () => {
            if (expectRenderedUserIDError) {
                it(`should show the error "${noUserIDError}"`, () => {
                    const error = screen.getByText(noUserIDError)
                    expect(error).toBeInTheDocument()
                    expect(error).toBeVisible()
                })
            } else {
                it(`should not show the error "${noUserIDError}`, () => {
                    const error = screen.queryByText(noUserIDError)
                    expect(error).not.toBeInTheDocument()
                })
            }
        })

        describe(`With expectRenderedCodePrefilled=${expectRenderedCodePrefilled}`, () => {
            if (typeof expectRenderedCodePrefilled == 'string') {
                it(`should show the ${codeLabelText} input with the value "${expectRenderedCodePrefilled}" prefilled`, async () => {
                    await waitFor(() => {
                        const codeInput = screen.getByLabelText(codeLabelText)
                        expect(codeInput).toHaveTextContent(expectRenderedCodePrefilled)
                    })
                })
            } else {
                it(`should not show the ${codeLabelText} input`, () => {
                    const codeInput = screen.queryByText(codeLabelText)
                    expect(codeInput).not.toBeInTheDocument()
                })
            }
        })

        describe(`With expectRenderedCodeError=${expectRenderedCodeError}`, () => {
            if (expectRenderedCodeError) {
                it(`should show the error "${noCodeError}"`, () => {
                    const error = screen.getByText(noCodeError)
                    expect(error).toBeInTheDocument()
                    expect(error).toBeVisible()
                })
            } else {
                it(`should not show the error "${expectRenderedCodeError}`, () => {
                    const error = screen.queryByText(noCodeError)
                    expect(error).not.toBeInTheDocument()
                })
            }
        })

        if (!directSubmit) {
            describe(`With click on submit`, () => {

            })
        }
    })
});