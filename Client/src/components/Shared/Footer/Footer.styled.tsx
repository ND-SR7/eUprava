import styled from "styled-components";

const FooterStyled = styled.footer`
  background-color: ${(props) => props.theme.colors.base};
  font-weight: ${(props) => props.theme.fontWeights.semiBold};
  font-style: italic;
  padding-bottom: ${(props) => props.theme.paddings.standard};
  position:absolute;
  left:0;
  bottom:0;
  right:0;
`;

export default FooterStyled;
