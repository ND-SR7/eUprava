import { useNavigate } from "react-router-dom";
import Button from "../components/Shared/Button/Button";
import HeadingStyled from "../components/Shared/Heading/Heading.styled";
import { useEffect, useState } from "react";
import { pingCourt } from "../services/PingService";
import toast from "react-hot-toast";
import Modal from "../components/Shared/Modal/Modal";
import { createHearingLegalEntity, createHearingPerson, createSuspension, createWarrant, getHearingsByJMBG, getSuspensionByJMBG, getWarrantsByJMBG } from "../services/CourtService";
import CourtHearing from "../models/Court/CourtHearing";
import SuspensionModel from "../models/Court/Suspension";
import Warrant from "../models/Court/Warrant";
import HearingList from "../components/Court/Hearing/HearingList";
import Suspension from "../components/Court/Suspension/Suspension";
import WarrantList from "../components/Court/Warrant/WarrantList";
import Form from "../components/Shared/Form/Form";

const CourtPage = () => {
  const navigate = useNavigate();

  const [modalVisible, setModalVisible] = useState(false);
  const [modalContent, setModalContent] = useState<any>(null);
  const [modalHeading, setModalHeading] = useState("");

  const hearingFormFields = [
    {attrName: "dateTime", label: "Date and time of hearing", type: "datetime-local"},
    {attrName: "reason", label: "Reason for hearing", type: "text"},
    {attrName: "jmbg", label: "Identification", type: "text"},
    {attrName: "type", label: "Person/Legal Entity", type: "radio", options: ["Person", "Entity"]}
  ];

  const suspensionFormFields = [
    {attrName: "from", label: "From date and time", type: "datetime-local"},
    {attrName: "to", label: "To date and time", type: "datetime-local"},
    {attrName: "person", label: "For person", type: "text"}
  ];

  const warrantFormFields = [
    {attrName: "trafficViolation", label: "Traffic violation ID", type: "text"},
    {attrName: "issuedFor", label: "Issued for person", type: "text"}
  ];
  
  useEffect(() => {
    const token = localStorage.getItem("token") || "";
    if (token === "") navigate("/");
  });

  const ping = () => {
    pingCourt().then(() => {
      toast.success("Good connection with court service");
    }).catch((error) => {
      toast.error("No connection to court service");
      console.error(error);
    });
  };

  const checkHearings = () => {
    getHearingsByJMBG().then((hearings: CourtHearing[]) => {
      if (hearings.length > 0) {
        setModalContent(<HearingList hearings={hearings} />);
        setModalHeading("Your hearings");
      } else {
        setModalContent("");
        setModalHeading("No court hearings for your account");
      }
      setModalVisible(true);
    }).catch(error => {
      toast.error("Failed to fetch hearings");
      console.error(error);
    });
  };

  const checkSuspensions = () => {
    getSuspensionByJMBG().then((suspension: SuspensionModel) => {
      if (suspension.id !== undefined) {
        setModalContent(<Suspension suspension={suspension} />);
        setModalHeading("Your suspension");
      } else {
        setModalContent("");
        setModalHeading("No suspension issued for your account");
      }
      setModalVisible(true);
    }).catch(error => {
      toast.error("Failed to fetch suspension");
      console.error(error);
    });
  };

  const checkWarrants = () => {
    getWarrantsByJMBG().then((warrants: Warrant[]) => {
      if (warrants !== null && warrants.length > 0) {
        setModalContent(<WarrantList warrants={warrants} />);
        setModalHeading("Your warrants");
      } else {
        setModalContent("");
        setModalHeading("No warrants issued for your account");
      }
      setModalVisible(true);
    }).catch(error => {
      toast.error("Failed to fetch warrants");
      console.error(error);
    });
  };

  const submitHearing = (formData: any) => {
    if (formData["type"] === "Person") {
      createHearingPerson({
        court: "64c13ab08edf48a008793cac",
        dateTime: formData["dateTime"] + ":00",
        reason: formData["reason"],
        person: formData["jmbg"]
      }).then(() => {
        toast.success("Successfully created a new hearing for person");
        closeModal();
      }).catch((error) => {
        toast.error("Failed to create a new hearing");
        console.error(error);
      });
    } else if (formData["type"] === "Entity") {
      createHearingLegalEntity({
        court: "64c13ab08edf48a008793cac",
        dateTime: formData["dateTime"] + ":00",
        reason: formData["reason"],
        legalEntity: formData["jmbg"]
      }).then(() => {
        toast.success("Successfully created a new hearing for legal entity");
        closeModal();
      }).catch((error) => {
        toast.error("Failed to create a new hearing");
        console.error(error);
      });
    }
  };

  const submitSuspension = (formData: any) => {
    createSuspension({
      from: formData["from"] + ":00",
      to: formData["to"] + ":00",
      person: formData["person"]
    }).then(() => {
      toast.success("Successfully created a new suspension");
      closeModal();
    }).catch((error) => {
      toast.error("Failed to create a new suspension");
      console.error(error);
    });
  };

  const submitWarrant = (formData: any) => {
    createWarrant({
      trafficViolation: formData["trafficViolation"],
      issuedFor: formData["issuedFor"]
    }).then(() => {
      toast.success("Successfully created a new warrant");
      closeModal();
    }).catch((error) => {
      toast.error("Failed to create a new warrant");
      console.error(error);
    });
  };

  const newHearing = () => {
    setModalHeading("Creating a new hearing");
    setModalContent(
      <Form 
        heading=""
        formFields={hearingFormFields}
        onSubmit={submitHearing} />
    );
    setModalVisible(true);
  };

  const newSuspension = () => {
    setModalHeading("Creating a new suspension");
    setModalContent(
      <Form 
        heading=""
        formFields={suspensionFormFields}
        onSubmit={submitSuspension} />
    );
    setModalVisible(true);
  };

  const newWarrant = () => {
    setModalHeading("Creating a new warrant");
    setModalContent(
      <Form 
        heading=""
        formFields={warrantFormFields}
        onSubmit={submitWarrant} />
    );
    setModalVisible(true);
  };

  const closeModal = () => {
    setModalVisible(false);
  };
  
  return (
    <>
      <HeadingStyled>Court</HeadingStyled>
      <br />
      <Button buttonType="button" label="Ping Service" onClick={() => ping()} />
      <br />
      <br />
      <Button buttonType="button" label="Check for hearings" onClick={() => checkHearings()} />
      <Button buttonType="button" label="Check for suspensions" onClick={() => checkSuspensions()} />
      <Button buttonType="button" label="Check for warrants" onClick={() => checkWarrants()} />
      <br />
      <br />
      <Button buttonType="button" label="Create hearing" onClick={() => newHearing()} />
      <Button buttonType="button" label="Create suspension" onClick={() => newSuspension()} />
      <Button buttonType="button" label="Create warrant" onClick={() => newWarrant()} />
      <br />
      <br />
      <Button buttonType="button" label="Go Back" onClick={() => window.history.back()}/>
      <br />

      <Modal
        heading={modalHeading}
        content={modalContent}
        isVisible={modalVisible}
        onClose={closeModal}
      />
    </>
  );
};

export default CourtPage;
