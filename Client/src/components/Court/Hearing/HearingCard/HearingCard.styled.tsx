import styled from "styled-components";

const HearingCardStyled = styled.section`
  background-color: ${(props) => props.theme.colors.base};
  border: ${(props) => props.theme.borders.standardBlack};
  margin: ${(props) => props.theme.margins.standard};
`;

export default HearingCardStyled;
