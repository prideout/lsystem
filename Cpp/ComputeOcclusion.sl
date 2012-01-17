class ComputeOcclusion(
                     string filename = "",
                     displaychannels = "",
                     coordsys = "";
                     color em = (1,0,1);
                     float samples = 256;
                     float Ka = 1;
                     float Kd =.5;
                     float Ks =.5;
                     float roughness =.1;
                     color specularcolor = 1;
)
{
  public void surface(output color Ci, Oi)
  {
    normal Nf = normalize(faceforward(N, I));
    vector V = vector(0) - normalize(I);
    normal Nn = normalize(Nf);

    float occ = occlusion(P, Nn, samples,
                          "maxdist", 100.0);

    if (filename != "")
        bake3d(filename, displaychannels, P, Nn,
                       "coordsystem", coordsys,
                       "_occlusion", occ);

    color c = em; /* (em + Ka*ambient() + Kd*diffuse(Nf)) +
                        specularcolor * Ks * specular(Nf, normalize(V), roughness); */

    Ci = (1 - occ) * Cs * c * Os;

    Oi = Os;
  }
}
