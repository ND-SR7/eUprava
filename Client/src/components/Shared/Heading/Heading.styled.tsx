import styled from "styled-components";

const HeadingStyled = styled.h1`
  color: ${(props) => props.theme.colors.base};
  font-weight: ${(props) => props.theme.fontWeights.bold};
`;

export default HeadingStyled;
